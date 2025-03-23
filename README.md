# unicorn

* it takes some time until a unicorn is produced, the request is blocked on requesters site and he need to wait.

* to improve the situation adjust the code, so that the requester is receiving a request-id, with this request-id he can poll and validate if unicorns are produced

* if the unicorn is produced it should be returned though using fifo principle

* adjust the code, so that every x seconds a new unicorn is produced at put to a store, which can be used to fulfill the requestQueue (LIFO Store)

* make sure, duplicate capabilities are not added to the unicorn

* improve the overall code

* if any requirements are not clear, compile meaningful assumptions

Requirements deduced from above:
1. Asynchronous Request Processing
2. Endpoints to Create, Check the status of request and as well as to retrieve unicorns using request-id
3. Unicorns will be returned in LIFO manner
4. Requests are processed in FIFO manner
5. Unicorn Should not have duplicate capabilities
6. Unicorn are produced after every x seconds and stored in Store