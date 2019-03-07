# dynamodb-copy
Copy one dynamodb to another and perform inserts and item edits along the way

##USAGE

make sure that aws-cli tools are installed and authenticated

build package and run ddb-copy with the correct flags

```console
cd dynamodb-copy
go build
cd cmd/ddb-copy
```


Example 1:

```console
ddb-copy -src payments-table-1 -dest payments-table-2 -aws-region us-east-1
```


##Command Line Options

-src  
source table
    
-dest  
  destination table
    
-aws-region  
aws region