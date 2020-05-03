# Bits of data

* Each file contains data of something specific. 
* Each file name must be in format `file-name.type.yaml`
* The `file-name` part might be anything that reflect it's data, it will act as ID of data
* The `type` part is a type of data and can be:
    * `generic`: Generic data, see `generic.sample.yaml`
* Files that ends with `.sample.yaml` or `.draft.yaml` will not be counted as data
* Files can be stored in sub-directories indefinitely
* Keep in mind that path to file is also a data ID, moving files between dirs, renaming dirs or files will also change their ID