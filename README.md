


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


## TOKEN
TODO
1. Define a Maker Interface with CreateToken and VerifyToken methods.
    - CreateToken(username string, duration time.Duration) (token, error)
    - VerifyToken(token) (*Payload, error)

2. Create a Payload struct that represents the token payload.
        -- ID(uuid)
        -- Username(string)
        -- IssuedAt(time.Time)
        -- ExpiredAt(time.TIme)
    2.1. Create a NewPayload function to return a new token payload for each specific username.
        -- NewPayload(username string, duration time.Duration) (*Payload, error)

3. Create a JWT implementation of the Maker
4. Write Unit test for the JWT implementation.
    -- Create a new JWTMaker by calling NewJWTMaker and pass in Random string and secretKey
    -- Generate a random username and duration, expiredAt and issuedAt fields
    -- Create a token
    -- Verify Token






## DOCKER
Create a Dockerfile to build the Golang application -- see Dockerfile.
    - The docker file uses multi-stage build to seperate building the original app and bundling the artifact.
    Using the multi-stage build method reduces the image size drastically.
* By default, the container will run on the default bridge network, this prevents the container from connecting with/to apps running on other containers.
To solve this, we need to connect both the App container and the Database container to the same app. That way, they can see each other and interact easily.

In this project, I connected the simplebank App container to the network created during initialization of the progressdb container (see docker-compose.yml)