## Log Collector in Go
## 
##
#### Problem statement

Write a TCP based log collection service where:
- Server:
    - Stores logs and buckets them as per the identifier sent by the client
    - Supports filesystem and S3 as backend storage
    - Supports log retrieval on text search
- Client:
    - Watches files for new logs
    - Forwards them to the server
    - Supports retry on failure

#### Design document
https://docs.google.com/document/d/1AZYdZ9kI26gsxA3NfVTGmfJKm_BdKypN-PExZQeK-Qc/edit?usp=sharing
