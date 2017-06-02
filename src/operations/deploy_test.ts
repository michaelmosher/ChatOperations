import 'mocha'
import { expect } from 'chai'

import * as Operations from './operations'
import * as Deploy from './deploy'

describe('Deploy.init()', () => {
    let subject = Deploy.init('michael')

    it('returns a new Deploy.Request', () => {
        expect(subject).to.be.an.instanceof(Deploy.Request)
    })
})

describe('Deploy.Request', () => {
    let subject : Deploy.Request
    let givenRequester = 'michael'

    beforeEach(() => {
        subject = Deploy.init(givenRequester)
    })

    describe('.action', () => {
        it('equals "deploy"', () => {
            expect(subject.action).to.equal('deploy')
        })
    })

    describe('.requester', () => {
        it('equals given Requester', () => {
            expect(subject.requester).to.equal(givenRequester)
        })
    })

    describe('.approved', () => {
        it('is initially undefined', () => {
            expect(subject.approved).to.equal(undefined)
        })
    })

    describe('.update()', () => {
        context('when updated with a Deploy.Server', () => {
            it('updates this.server', () => {
                let server = new Deploy.Server('dev', '123.45.6.789')
                subject.update(server)
                expect(subject.server).to.equal(server)
            })
        })

        context('when updated with a Operations.Responder', () => {
            it('updates this.responder', () => {
                let responder = new Operations.Responder('admin')
                subject.update(responder)
                expect(subject.responder).to.equal('admin')
            })
        })

        context('when updated with a Operations.Approved', () => {
            it('updates this.approval', () => {
                let responder = new Operations.Approval(true)
                subject.update(responder)
                expect(subject.approved).to.equal(true)
            })
        })

        context('when updated with a Operations.Success', () => {
            it('updates this.succeeded', () => {
                let responder = new Operations.Success(true)
                subject.update(responder)
                expect(subject.succeeded).to.equal(true)
            })
        })

        context('when updated with a Operations.ResponseURL', () => {
            it('updates this.response_url', () => {
                let responder = new Operations.ResponseURL('webhook.slack.com')
                subject.update(responder)
                expect(subject.response_url).to.equal('webhook.slack.com')
            })
        })
    })

    describe('.isReady', () => {
        context('when .server is undefined', () => {
            it('returns false', () => {
                expect(subject.server).to.be.undefined
                expect(subject.isReady()).to.be.false
            })
        })

        context('when .response_url is undefined', () => {
            it('returns false', () => {
                expect(subject.response_url).to.be.undefined
                expect(subject.isReady()).to.be.false
            })
        })

        context('when both are defined', () => {
            beforeEach(() => {
                subject.update(new Operations.ResponseURL('any.legal.url'))
                subject.update(new Deploy.Server('test', 'any.test.url'))
            })

            it('returns true)', () => {
                expect(subject.isReady()).to.be.true
            })
        })
    })

    describe('.next', () => {
        it('returns a list of servers', () => {
            let options = subject.next()
            expect(options).to.eql([new Deploy.Server('dev', '123.4.56.789')])
        })
    })

    describe('.summary', () => {
        context('when .server is undefined', () => {
            it('returns a partial request summary message', () => {
                let message = subject.summary()
                expect(message).to.equal('michael has requested a deploy.')
            })
        })

        context('when .server is defined', () => {
            it('returns a complete request summary message', () => {
                subject.update(new Deploy.Server('dev', '123'))
                let message = subject.summary()
                expect(message).to.equal('michael has requested a deploy on dev.')
            })
        })
    })
})
