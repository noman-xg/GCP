# s204
## _Implementation of  ETL pipeline on GCP Infrastructure_



![Build Status](https://user-images.githubusercontent.com/95742163/150748975-dc63935f-5644-4edf-8c74-2a71aa76846d.png)

In this cloud specialization module we'll be implementing an ETL (Extract,Transform,Load) pipeline through a Google cloud function written in GOLANG. We'll be using Terraform as IAC tool to provision and use following GCP resources for our ETL pipeline. 

- Google Cloud Storage.
- Google CLoud Function.
- Google Cloud Scheduler.
- GCP BigQuery

---
## Features

- Extract 100 csv records per function invocation.
- Transform records which meet the transformation criteria.
- Insert the records into the BigQuery Table. 
- Enable periodic function invocation using Google Cloud Scheduler.
---
## Step-by-Step Guide

1 - Clone the xldp repo to your system and navigate to the directory relevant to this module.

```sh
git clone https://github.com/X-CBG/xldp.git
cd xldp/cloud_specializations/s204/noman
```

2 - Download the csv file from [this](https://www.stats.govt.nz/assets/Uploads/Effects-of-COVID-19-on-trade/Effects-of-COVID-19-on-trade-At-15-December-2021-provisional/Download-data/effects-of-covid-19-on-trade-at-15-december-2021-provisional.csv) link. 

3 - Navigate to the CloudFunction directory and create a zip archive (_`ETL.zip`_) of the cloud function source code files.
```sh
cd CloudFunction/
zip ETL.zip main.go go.mod && cd ..
```

4 - Navigate to the Terraform folder. initialize terraform in the directory.
```sh
cd Terraform && terraform init
```
5 - Open the directory in any editor for example VScode and update the values of _`sourcepath`_ and _`csv_local_path`_  
    variables in `variables.tf` according to the absolute path of the csv file and ETL.zip on your system. Futher, please update the value of `project` in `main.tf` to reflect the ProjectID of your project in your GCP account.  

6 - Authenticate your gcloud cli with GCP. 
```sh
gcloud auth login 
#you will be prompted with a window in your default browser, check the boxes as per requirements and your gcloud CLI will be authenticated with GCP automatically.
```
7 - Apply the terraform configuration to setup the infrastructure resources
```sh
terraform apply --auto-approve
```
---
## GCP Resources

Login into your gcp account and verify that the following resources with the mentioned names have been created inside your project.

| Resource | Name |
| ------ | ------ |
| Cloud Storage Bucket | `for_s204_xgrid` |
| Cloud Storage Object | `covid-csv-s204.csv` |
| Cloud Function | `etl_function` |
| Cloud Scheduler| `invoke_ETL` |
| BigQuery Dataset | `covidDataset` |
| BigQuery Table | `covid-table` |

---
## Verify
You may verify the functionality by manually triggering the cloud fucntion http endpoint from your browser.

[_http://us-east1-`your-projectID`.cloudfunctions.net/etl_function/?a=`your-projectID`&&b=covidDataset&&c=covid-table_](_http://us-east1-`your-projectID`.cloudfunctions.net/etl_function/?a=`your-projectID`&&b=covidDataset&&c=covid-table_)

** _Note: Please replace "_your-projectID_" with the projectID of your project in your GCP accout._

You can now verify the functionality of the ETL pipeline by previewing the BigQuery table in the dataset which will now be including the first 100 records from the csv after the first invocation.

---



