'use strict';

// import everything!
const AWS = require('aws-sdk')
const Request = require('request')

const RequestDB = require('./lib/database/requests')
const Slack = require('./lib/notification/slack')
const Application = require('./lib/application/requestManager')

// create database
const awsDDB = new AWS.DynamoDB.DocumentClient({
    region: 'us-east-2',
    endpoint: 'http://localhost:8000'
})
const requestRepo = new RequestDB.Repo(awsDDB)

// create slack notifier
const slackNotifier = new Slack.Notifier(Request, 'webhook.slack.com') // TODO

// create application
const requestManager = new Application.Manager(requestRepo, slackNotifier)

// create webserver
