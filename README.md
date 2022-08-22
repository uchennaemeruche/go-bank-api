


There is an issue with the mockgen package -- such that it doesn't get added to the vendor list automatically.
Also, it gets removed anytime you run ``` go mod tidy ```
As a workaround, In this project, I added an empty import for ``` _ "github.com/golang/mock/mockgen/model" ``` and then ran ``` go mod tidy```