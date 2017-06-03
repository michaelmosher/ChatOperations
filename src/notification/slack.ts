import * as Operations from '../operations/operations'
import * as Templates from './templates'

declare type callback = (err: Error, req: any, body: any) => void

interface WebClient {
    post(options: any, cb?: callback): void
}

export class Notifier {
    constructor(public http: WebClient, readonly hostingWebhook: string) { }

    reportError(url: string, err: Error) : void {
        let body = Templates.errorReport
        body.attachments[0].text = err.message
        return this.http.post({ url: url, body: body })
    }

    handleErrors(url: string) : callback {
        return (err: Error, req: any, body: any) => {
            if (err) {
                return this.reportError(url, err)
            }
            if (req.statusCode !== 200) {
                return this.reportError(url, body.error)
            }
        }
    }

    requestSubmitted(r: Operations.Request) : void {
        let body = Templates.submitRequest
        body.text = r.summary()
        body.attachments[0].callback_id = r.callback_id

        this.http.post({
            url: this.hostingWebhook, 
            body: body
        }, this.handleErrors(r.response_url))
    }

    requestAnswered(r: Operations.Request, approved: boolean) : void {
        let body = {
            text: '',
            replace_original: false
        }

        body.text = (approved)
            ? `:white_check_mark: ${r.responder} approved your request.`
            : `:x: ${r.responder} denied your request.`

        this.http.post({
            url: r.response_url,
            body: body
        }, this.handleErrors(this.hostingWebhook))
    }

    requestApproved(r: Operations.Request) : void {
        this.requestAnswered(r, true)
    }

    requestDenied(r: Operations.Request) : void {
        this.requestAnswered(r, false)
    }
}