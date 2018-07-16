# smpl-mesg-service
mesg-service-test as Proof of Concept to create common platform to integrate various technologies seamlessly with user application.

Currently, this service supports receiving POST requests from HTTP clients, raising corresponding event with mesg-core and responding to client with relevant response from the upstream services.

## Assumptions
mesg-service-test sends event to mesg-core with the body of the response made to the upstream stream

## Contracts
The /batchexecute is POST method only endpoint which receives request body in the following format:
`[
    {URL:...
    Body:...},
    {URL:...
    Body:...}
]`. Where URL is the URL of the upstream endpoint and Body is the payload for the request.

And following is the Response received once all the requests in the batch execution has executed, either normally or exceptionally:
`{
    URL:{
        StatusCode:...
        Body:...
    },
    URL:{
        StatusCode:...
        Body:...
    }
}` .Where each url in the batch request is represented by `URL` and corresponding response body by `Body` and status code by `StatusCode`