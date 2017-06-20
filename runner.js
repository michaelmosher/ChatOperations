const AWS = require('aws-sdk')

const Request = require('./lib/database/request')
const deploy = require('./lib/operations/deploy')

var docClient = new AWS.DynamoDB.DocumentClient({
    region: 'us-east-2',
    endpoint: 'http://localhost:8000'
})

var requestRepo = new Request.Repo(docClient)
var callback_id = 'deploy_michael_47486376';

var req = deploy.init('michael')

console.log(requestRepo.putParamBuilder(req))

