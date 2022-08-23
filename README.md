


There is an issue with the mockgen package -- such that it doesn't get added to the vendor list automatically.
Also, it gets removed anytime you run ``` go mod tidy ```
As a workaround, In this project, I added an empty import for ``` _ "github.com/golang/mock/mockgen/model" ``` and then ran ``` go mod tidy```


## TODO:
- Design Load Test for the entire API (How the API behaves under an expected and high load or requests)
    In Load Testing, we should stimulate actual user load on the API.
- Design Stress Test for the API (How the API will behaves using a load beyond the expected maximum - such as a DDoS attack, Slashdot effect, or other scenarios)
    Stress Tests determines the stability and robustness of the system. How the system behaves under extreme loads and how it recovers from failures.
- Brutforce Test the API
- Add these tests to CI/CD process

https://www.guru99.com/performance-vs-load-vs-stress-testing.html
https://www.artillery.io/pricing
https://www.loadmill.com/
https://jmeter.apache.org/