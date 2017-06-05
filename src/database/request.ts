import * as AWS from 'aws-sdk'
import * as Deploy from '../operations/deploy'
import * as Operations from '../operations/operations'

interface awsDynamoDBDocClient {
    get(params: AWS.DynamoDB.DocumentClient.GetItemInput, 
        callback?: ((err: AWS.AWSError, data: AWS.DynamoDB.DocumentClient.GetItemOutput) => void) | undefined): AWS.Request<AWS.DynamoDB.DocumentClient.GetItemOutput, AWS.AWSError>,
    put(params: AWS.DynamoDB.DocumentClient.PutItemInput, 
        callback?: (err: AWS.AWSError, data: AWS.DynamoDB.DocumentClient.PutItemOutput) => void): AWS.Request<AWS.DynamoDB.DocumentClient.PutItemOutput, AWS.AWSError>
}

export class Repo {
    constructor(readonly dynamodb: awsDynamoDBDocClient) { }

    protected getParamBuilder(callback_id: string): AWS.DynamoDB.GetItemInput {
        return { 
            TableName: 'Requests',
            Key: { 'callback_id': callback_id }
        }
    }

    protected async getRequestItem(callback_id: string): Promise<any> {
        return new Promise((resolve, reject) => {
            let params = this.getParamBuilder(callback_id)

            this.dynamodb.get(params, (error, data) => {
                if (error) {
                    return reject(error)
                }

                return resolve(data['Item'])
            })
        })
    }

    protected parseRequestItem(item: {}): Operations.Request {
        // TODO check item.action
        let deployRequest = Deploy.init('')
        return Object.assign(deployRequest, item)
    }

    async findById(callback_id: string): Promise<Operations.Request> {
        let requestItem = await this.getRequestItem(callback_id)
        return this.parseRequestItem(requestItem)
    }

    protected putParamBuilder(r: Operations.Request): AWS.DynamoDB.PutItemInput {
        let param: {TableName: string, Item: any} = {
            TableName: 'Requests', 
            Item: {}
        }

        for (let key in r) {
            if (typeof r[key] !== 'function') {
                param.Item[key] = r[key]
            }
        }

        return param
    }

    protected async putRequestItem(r: Operations.Request): Promise<any> {
        return new Promise((resolve, reject) => {
            let params = this.putParamBuilder(r)

            this.dynamodb.put(params, (error, data) => {
                if (error) {
                    return reject(error)
                }

                return resolve()
            })
        })
    }

    async store(r: Operations.Request): Promise<Operations.Request> {
        return this.putRequestItem(r)
        .then(() => {
            return r
        })
    }
}