import * as Operations from '../operations/operations'
import * as Templates from './templates'

declare type callback = (err: Error, req: any, body: any) => void

interface WebClient {
    post(options: any, cb?: callback): void
}

export class Notifier {
    constructor(public http: WebClient, readonly hostingWebhook: string) { }

    async promPostRequest(options: any): Promise<any> {
        return new Promise((resolve, reject) => {
            this.http.post(options, (err, req, body) => {
                if (err) {
                    return reject(err)
                }
                if (req.statusCode !== 200) {
                    return reject(body.error)
                }
                return resolve(body)
            })
        })
    }

    reportError(url: string, err: Error) : Promise<any> {
        let body = Templates.errorReport
        body.attachments[0].text = err.message
        return this.promPostRequest({ url: url, body: body })
    }

    requestSubmitted(r: Operations.Request) : Promise<any> {
        let body = Templates.submitRequest
        body.text = r.summary()
        body.attachments[0].callback_id = r.callback_id

        return this.promPostRequest({
            url: this.hostingWebhook, 
            body: body
        })
        .catch(error => {
            this.reportError(r.response_url, error)
        })
    }

    requestAnswered(r: Operations.Request, approved: boolean) : Promise<any> {
        let body = {
            text: '',
            replace_original: false
        }

        body.text = (approved)
            ? `:white_check_mark: ${r.responder} approved your request.`
            : `:x: ${r.responder} denied your request.`

        return this.promPostRequest({
            url: r.response_url, 
            body: body
        })
        .catch(error => {
            this.reportError(this.hostingWebhook, error)
        })
    }

    requestApproved(r: Operations.Request) : void {
        this.requestAnswered(r, true)
    }

    requestDenied(r: Operations.Request) : void {
        this.requestAnswered(r, false)
    }
}