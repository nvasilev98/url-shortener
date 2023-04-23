## Decisions
1. Database:

   Expecting a high amount of traffic and requiring a large storage size, while not anticipating any joins, suggests that a NoSQL database would be suitable. Firestore has been selected for this particular use case.

2. Encoding algorithm

   There is a limitation to use only Latin letters and digits, which means there are 62 different symbols - [0–9][a-z][A-Z]. As we want to support short URLs with max length - 5 generating 62^5 unique short urls equals approximately to ~915 million. This means if a unique number is provided everytime, base62 encoding will be sufficient.

3. Distributed Counter

   As NoSQL DB Firestore does not have built-in auto-increment operator, so we need a way to guarantee that a unique number is used every time generating a new short URL. For this purpose, a distributed counter is implemented in Firestore which is a collection of shards. Each shard has its own count field and when the counter is incremented, a random shard is used.

    Reference: https://firebase.google.com/docs/firestore/solutions/counters

4. Router

    For the router is used - Gin. A widely used web framework that provides performance, readability. Gin is MIT licensed and actively maintained.

    Reference: https://github.com/gin-gonic/gin

5. Unit tests

    For the unit tests is used - Ginkgo & Gomega. A widely used unit test frameworks which provide readability, expressive and structured unit tests. Both Ginkgo & Gomega are MIT licensed as well as actively maintained.

    For mocking in unit tests is used - Gomock. A mocking framework that provides readability in the unit tests.

    Reference: 
    1. https://onsi.github.io/ginkgo/
    2. https://onsi.github.io/gomega/
    3. https://github.com/golang/mock

6. Project layout

    For the project layout some good practices have been applied from: https://github.com/golang-standards/project-layout