package etlpipeline

import (
	"context"
	"encoding/csv"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/storage"
)

//These variables being in global (instance-wide) scope will be used for state-preservation across the function executions.
var startingIndex int
var endingIndex int

//Struct for the CSV records to be loaded into the BigQuery Table.
type csvRecord struct {
	Direction      string
	Year           string
	Date           string
	Weekday        string
	Country        string
	Commodity      string
	Transport_Mode string
	Measure        string
	Value          string
	Cumulative     string
	code_index     int
}

//init runs during package initialization. So, this will only run during the instance's cold start and is used to initialize the startingIndex variable.
func init() {
	startingIndex = 1
}

func ETL(w http.ResponseWriter, r *http.Request) {

	/*ETL function is the entryPoint function of our cloud . It will handle the HTTP request
	made to the exposed endpoint and then we'll extract and use the parameters passed in
	the HTTP GET request as ProjectID,DatasetID and TableID. It will help ensure the generic
	nature of our cloud fucntion.
	*/

	params := r.URL.Query()
	projectID := params["a"][0]
	datasetID := params["b"][0]
	tableID := params["c"][0]

	fmt.Fprint(w, html.EscapeString(projectID))
	fmt.Fprint(w, html.EscapeString(datasetID))
	fmt.Fprint(w, html.EscapeString(tableID))

	ExtTrans(projectID, datasetID, tableID)
}

func ExtTrans(projectID, datasetID, tableID string) error {

	/* ExtTrans function performs two basic fucntions of the Extract-Transform-Load Pipeline.
	1- Reads the CSV data from the Google Storage bucket into the program memory.
	2- Checks and  transforms every csv record if the Weekday in Record has Sunday as it's value.
	3- Passes over the csv record to the loadInBigQuery which subsequently puts the record in the BigQuery table.
	*/

	//Google Cloud Storage bucket details.
	const BUCKET = "for_s204_xgrid"
	const OBJECT = "covid-csv-s204.csv"

	//Creates context and client for GCS.
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Printf("stage 1.%v", err)
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	fmt.Println("Stage 1. CLient and context created.")

	//Creates a reader object to read CSV file from the Google Cloud Storage bucket.
	reader, err := client.Bucket(BUCKET).Object(OBJECT).NewReader(ctx)
	if err != nil {
		fmt.Println(err)
	}

	rec := csv.NewReader(reader)

	//ReadAll reads the complete CSV into data.
	data, err := rec.ReadAll()
	if err == io.EOF {
		fmt.Printf(" Error while resding the file: %v", err)
	}
	fmt.Println("Stage 2. Read from Storage done.")

	//The global varibales: startingIndex and endingIndex will govern the range of csv records that we will be loading into the BigQuery per function invocation.
	endingIndex = startingIndex + 100 // since we want to 100 records/execution.

	// for loop to iterate over the nested string array (data) containing the read CSV records.
	for _, record := range data[startingIndex:endingIndex] {

		if err != nil {
			log.Fatal(err)
		}

		// Transform logic as per requirement.
		if record[3] == "Sunday" {
			record[8] = "0"
			record[9] = "0"
			fmt.Println(record)

			// calls the loadInBigQuery to insert the transformed record into the BigQuery table.
			err = laodInBigQuery(record, projectID, datasetID, tableID)
			if err != nil {
				fmt.Printf("Error while loading record %v : %v: ", record, err)
			}
		} else {
			err = laodInBigQuery(record, projectID, datasetID, tableID) //calling loadInBugQuery for records which do not qualify to be transformed as per condition.
			if err != nil {
				fmt.Printf("Error while loading record %v : %v: ", record, err)
			}
		}
	}

	//updating the value of startingIndex as current endingIndex so that it can be preserved and subsequently be used as startingIndex value for next function invocation.
	startingIndex = endingIndex
	fmt.Printf("the starting index for next execution is: %v", startingIndex)
	return nil
}

//loadinBigQuery inserts the data record by record into a BigQuery table with pre-defined schema.
func laodInBigQuery(record []string, projectID, datasetID, tableID string) error {

	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		fmt.Printf(" error in loadintoBigQUery: %v", err)
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()
	//creates a variable of the type csvRecord (struct decalared in global scope) and populates with the values as in the record that has been passed to loadinBigQuery.
	csvRecords := []*csvRecord{
		{Direction: record[0], Year: record[1], Date: record[2], Weekday: record[3], Country: record[4],
			Commodity: record[5], Transport_Mode: record[6], Measure: record[7], Value: record[8],
			Cumulative: record[9]}}

	//Creates an inserter object
	inserter := client.Dataset(datasetID).Table(tableID).Inserter()
	//calling .put on inserter to insert record into the BIgQuery table
	if err := inserter.Put(ctx, csvRecords); err != nil {
		fmt.Printf("Error in .put : %v", err)
	}
	return nil
}
