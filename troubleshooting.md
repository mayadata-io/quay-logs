# Quay-logs - Troubleshooting

- You can enable Debug flags while executing `go run cmd/main.go --debug=true` to get the insights of errors.
- If you are unable to run `main.go` ensure that you have entered the right access token and namespace.
  -- To get the access token contact the namespace owner or read [quay_faq.md](https://github.com/mayadata-io/quay-logs/blob/master/quay_faq.md) in the repo.
- Read the `how_to_contribute.md` to get the detailed insights of the code.
- If any new code is added then if possible make a `Debug` if block which will only run when debug flag is enabled.

## Some points to note

- The filenames in `./logs` are in the format `MonDDYYYYDateFormat = "Jan-02-2006"` i.e example `Aug-06-2020-09_13_10-2.json`(where 09_13_10 is the time of file creation).
  -- each of this json has around 20 `Logs` at max. So there are multiple pages. That's why you can see the above example has `-2`at the end, which denotes the pagenumber/token/index.

---

These codes are being run on a regular basis to gather the data/logs so that we can send in those data to prometheus and then connect it with grafana to get meaningful graphs.

---

For quay-logs FAQs refer [quay_faq.md](https://github.com/mayadata-io/quay-logs/blob/master/quay_faq.md) file in the repo.

Contributing information is available at `how_to_contribute.md` file.
