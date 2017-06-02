import 'mocha'
import { expect } from 'chai'

import * as Slack from './slack'
import * as Operations from '../operations/operations'
import * as Deploy from '../operations/deploy'

describe('Slack.Notifier', () => {
    type callback = (err: Error, req: any, body: any) => void
    let dummyWebClient = {
        post: function(options: any, cb: callback) : void {

        }
    }
    let slackNotifier: Slack.Notifier
    let request: Operations.Request

    beforeEach(() => {
        slackNotifier = new Slack.Notifier(dummyWebClient, 'hosting.slack.com')
        request = Deploy.init('michael')
    })

    describe('.reportError', () => {
        beforeEach(() => {
            slackNotifier.http.post = function(options: any, cb: callback) : void {
                expect(options.url).to.equal('dummy.com')
                expect(options.body.attachments[0].text).to.equal('problems!')
            }
        })

        it('sends error notification to the given url', () => {
            slackNotifier.reportError('dummy.com', new Error('problems!'))
        })
    })

    describe('.handleErrors', () => {
        it('returns a function to handle http request errors', () => {
            let handler: callback = slackNotifier.handleErrors('some.url.com')
            expect(handler).to.exist
        })
    })

    describe('.requestSubmitted', () => {
        beforeEach(() => {
            slackNotifier.http.post = function(options: any, cb: callback) : void {
                expect(options.url).to.equal('hosting.slack.com')
                expect(options.body.text).to.include('michael has requested a deploy on dev')
            }
        })
        it('sends notification to the hostingWebhook', () => {
            request.update(new Deploy.Server('dev', 'dev.com'))
            slackNotifier.requestSubmitted(request)
        })
    })

    describe('.requestApproved', () => {
        beforeEach(() => {
            slackNotifier.http.post = function(options: any, cb: callback) : void {
                expect(options.url).to.equal('webhook.slack.com')
                expect(options.body.text).to.include('admin approved your request')
            }
        })
        it('sends approval notification to a Request\'s response_url', () => {
            request.update(new Deploy.Server('dev', 'dev.com'))
            request.update(new Operations.Responder('admin'))
            request.update(new Operations.ResponseURL('webhook.slack.com'))
            slackNotifier.requestApproved(request)
        })
    })

    describe('.requestDenied', () => {
        beforeEach(() => {
            slackNotifier.http.post = function(options: any, cb: callback) : void {
                expect(options.url).to.equal('webhook.slack.com')
                expect(options.body.text).to.include('admin denied your request')
            }
        })
        it('sends rejection notification to a Request\'s response_url', () => {
            request.update(new Deploy.Server('dev', 'dev.com'))
            request.update(new Operations.Responder('admin'))
            request.update(new Operations.ResponseURL('webhook.slack.com'))
            slackNotifier.requestDenied(request)
        })
    })
})