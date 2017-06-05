import 'mocha'
import { expect } from 'chai'
import * as AWS from 'aws-sdk'

import * as Request from './request'

import * as Deploy from '../operations/deploy'
import * as Operations from '../operations/operations'

describe('Request.Repo', () => {
    let subject: Request.Repo
    let testRequest: Operations.Request
    let testDynamoDb = new AWS.DynamoDB.DocumentClient({
        region: 'us-east-2',
        endpoint: 'http://localhost:8000'
    })
    
    before(() => {
        testRequest = Deploy.init('test.user')
    })

    beforeEach(() => {
        subject = new Request.Repo(testDynamoDb)
    })

    describe('.store', () => {
        it('stores a request in dynamodb', (done) => {
            subject.store(testRequest)
            .then(req => {
                expect(req).to.eql(testRequest)
                done()
            })
        })
    })

    describe('.findById', () => {
        it('reads a request from dynamodb', (done) => {
            subject.findById(testRequest.callback_id)
            .then(req => {
                expect(req).to.eql(testRequest)
                done()
            })
        })
    })
})