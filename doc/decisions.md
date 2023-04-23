## Decisions
1. Database:

    We expect high amount of traffic and we will need huge storage size and we do not need any joins so we can select a NoSQL database. In this case we are going to use Firestore.

2. Encoding algorithm

   We have a limitation to use only Latin letters and digits, which means we have 62 different symbols - [0–9][a-z][A-Z]. As we want to support short URLs with max length - 5 we can generate 62^5 unique which is approximately equal to ~915 million. This means if we can guarantee a unique number everytime we can simply use base62 encoding for this purpose.

3. Distributed Counter

    As NoSQL DB Firestore does not have built-in auto increment operator, so we need a way to guarantee that we use a unique number every time we generate a new short URL. For this purpose a distributed counter is implemented in Firestore which is a collection of shards. Each shard has its own count field and when we increment the counter a random shard is used.

    Reference: https://firebase.google.com/docs/firestore/solutions/counters

4. Router

    For the router we have used - Gin. A widely used web framework that provides performance, readability. Gin is MIT licensed and actively maintained.

    Reference: https://github.com/gin-gonic/gin

5. Unit tests

    For the unit tests we have used - Ginkgo & Gomega. A widely used unit test frameworks which provide readability,expressive and structured unit tests. Both Ginkgo & Gomega are MIT licensed as well as actively maintained.

    Reference: 
    1. https://onsi.github.io/ginkgo/
    2. https://onsi.github.io/gomega/

