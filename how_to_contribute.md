# Quay-logs - Contributing/Docs
### cmd/main.go file
>The possible list of command line arguments/flags which can be given to the `go run cmd/main.go` command are listed inside the var( ) method.
 - The main function has the following logic
 -- It makes the required directories.
 -- It lists all the repos in the sorted order of popularity in the namespace into `repolist`.
 -- It iterated through each of the repos and download its Logs and stores the in different files.
 
- **flag.Parse( )** : parses the flags. It must be called before using any of the flags.
- **mkdirAll( )**: makes all the required directories/folders from the arguments  by formatting them with -p flags.
-- Here we create 1 directory/folders
**/logs** : for storing the repos in order of popularity

- We create a `NewLister` (refer `list.go`) and set `IsWriteToFile` false because we don't want to store the data in files.
- We call the `ListReposAndWriteToFileOptionally( )`function to get all the repolist in the namespace as a JSON. It returns all the repos in sorted order of popularity.
- --
### list.go file
- **NewLister( )** - It creates a new folder by the mkdir command using the arguments passed to it. It returns a type of new Listable.

- **ListReposAndWriteToFileOptionally( )** -actually calls `ListReposByPopularityAndWriteToFileOptionally( )`function. 
- **ListReposByPopularityAndWriteToFileOptionally( )** - It calls the `RequestReposForPageToken( )` which returns all repos name in order of popularity.
-- Right now we don't have 100 repos that's why all the data are in one page. Thus some codes are commented.

- **RequestReposForPageToken( )** - Creates a HTTPRequest with some query parameters and invokes it. 
-- Since `IsWriteToFile` is false so it **doesn't** call`WriteToFile`) and the JSON is unmarshaled and returned. 

- **WriteToFile( )** - Writes the content of response body into passed filename with file mode 0644.
- --
### types.go file
- It has all the required structures which can be used to extract the data from quay(JSON).
- The structures are in the order of
-- PopularList -> Popular
-- LogList -> Log -> Metadata -> ResolvedIP

---
### logs.go file
- **NewLogger( )** -  It creates a new folder by the mkdir command using the arguments passed to it for each of the repos. Example: `./logs/namespace/reponame/` It returns a type of new Listable.
- **Log( )** - It calls `RequestLogsForPageToken( )` to get the logs from the Quay API. It stores them in separate files by calling `WriteToFile` internally. 
--Here next page is available since the API returns 20 `logs` at once. So each files can contain at max 20 `logs`.
- **RequestLogsForPageToken( )** -  Creates a HTTPRequest with some query parameters and invokes it. 
-- Since `IsWriteToFile` is true here so it calls `WriteToFile`) and the JSON is unmarshaled and returned. 

- **WriteToFile( )** -  Writes the content of response body into passed filename with file mode 0644. It stores the logs into `./logs/namespace/reponame/filename.json`
>**Why we have created `/logs` folder🤔?**
>Since the data in `./logs` folder we are creating, is used by the [growth-metrics](https://github.com/mayadata-io/growth-metrics) repo to extract the `logs` and then `sanitise` the logs to produce meaningful data. Which then can be sent to Prometheus and Grafana to display it in graphs.
---
### http_client.go file
- It has all the necessary functions to make a **HTTPRequest** and **Invoke** it. 
- --
For quay-logs FAQs refer [quay_faq.md](https://github.com/mayadata-io/quay-logs/blob/master/quay_faq.md) file in the repo.
Troubleshooting information is available at `troubleshooting.md` file.