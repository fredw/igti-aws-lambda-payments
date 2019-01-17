Payments process using AWS Lambda 
============================================================

The digital payment made through the internet using credit cards or bank transactions has become essential for success in a business, providing business and consumer advantages by making payment easily and securely. 
To provide this functionality, there are a few different market models, one of which is integration through a payment gateway, a transparent model for the consumer. 

The complexity of this integration is high, which can have serious consequences and directly impact the company's profits. 
To ensure that these problems are analyzed in the right way, it is necessary to control this processing, managing each case in a particular way. 
This project demonstrates a hypothetical solution of integration with a payment gateway using the **Serverless** architecture through **AWS Lambda**. 
The communication with the gateway may lead to possible failures and possible processing problems, it is crucial to guarantee, control, reprocess or restrict processing in case of critical failures, making it necessary to store the status of the transaction, made through **Amazon Simple Queue Service**.


Development
------------------------------------------------------------

### Environment

These are the available and used environment variables that are used inside the **AWS Lambda** function:

* `LOG_LEVEL`: the log level. Possible values: `INFO`, `DEBUG`, `WARNING`, etc. (default: `INFO`);
* `SQS_QUEUE_URL`: the SQS Queue URL to consume the payment messages (`required`); 
* `SQS_DLQ_QUEUE_URL`: the Dead Letter Queue SQS Queue URL, used to move the messages that were processed and have critical errors (`required`); 
* `SQS_MAX_NUMBER_OF_MESSAGES`: the maximum number of messages that will be read for each execution of the function (`required` and the default value is `1`); 
* `PROVIDER_EXAMPLE_REQUEST_URI`: the URL used to integrate the payments with the `Example` provider. As this project uses an hypothetical integration situation, we use this `Example` url with mocked results; 

### Commands

To run unit tests:
```bash
make test
```

To build the project:
```bash
make build
```

To updated the project on AWS Lambda*:
```bash
make update
```

To invoke the AWS Lambda function*:
```bash
make invoke
```

**Test**: read all messages from `payments.fifo` queue on SQS:
```bash
make sqs_receive_messages
```

**Test**: create a test message on `payments.fifo` queue on SQS:
```bash
make sqs_create_test_message
```

**Test**: purge all messages from `payments.fifo` queue on SQS:
```bash
make sqs_purge_queue
```

*You need setup your `aws cli` credentials to be able to use that and have right permissions to be able to do that.


Credits
------------------------------------------------------------
* Frederico Wuerges Becker <fred.wuerges@gmail.com>
 