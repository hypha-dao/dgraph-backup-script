# Dgraph Backup Script

Creates a cron export job for each of the specified export jobs in the configuration file.

To run the dgraph backup script with the helper script:

`./start-local-env.sh`

Or:

`go run . ./config.yml`

Look at the **config.yml** file for an example configuration file, some important configuration parameters are:

- export-destination: The endpoint where the backups will be stored, currently minio as bridge to store the backups in google storage
- export-access-key: The access key for the specified endpoint
- export-secret-key: The password or secret for the specified endpoint
- export-jobs: The configuration for the export jobs to run
  - name: name of the job, used as dir name when storing the backup
  - gql-admin-url: dgraph's admin endpoint
  - schedule: cron time format to specify when the job should run 


To run dgraph backup script by building the dgraph backup script image based on the current state of the project:

`./start-dgraph-backup-script build`

To run dgraph backup script by pulling the latest dgraph backup script image from docker hub:

`./start-dgraph-backup-script image`

