# quay-logs

## Step 1/
```sh
#  This will create the binary named main
go build cmd/main.go
```

## Step 2/
```sh
# Invoke the binary with quay auth token & repository
#
# - NOTE: Quay auth token need to be provided
# - NOTE: Namespace that hosts the images should be provided
#
# Following activities are handled by this binary:
#
# - Downloads latest quay images with popularity ranks
# - Downloads latest quay logs for each image
# - Sanitises these logs by removing duplicate entries
# - Enhances the sanitised logs by adding ip domain information
./main --quay-auth-token=<auth token> --quay-namespace=openebs
```

**Note:** quay-auth-token should have scope of `Administer Repositories`.

## Folder details
- **logs/** has actions on each image categorized by dates

## Few quay.io APIs w.r.t openebs
- https://quay.io/api/v1/repository?popularity=true&namespace=openebs
- https://quay.io/api/v1/repository/openebs/provisioner-localpv/logs

## Source code details
- Refer to **cmd/main.go** for various arguments that can be provided to this binary
- **list.go** has the logic to download current popularity/ranking logs of quay namespace
- **logs.go** has the logic to download quay image logs based on a date range
- **types.go** has quay API schema coded as go structure
