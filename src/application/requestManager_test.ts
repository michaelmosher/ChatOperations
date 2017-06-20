import 'mocha'
import { expect } from 'chai'

import * as Request from './requestManager'
import * as Deploy from '../operations/deploy'
import * as Operations from '../operations/operations'

describe('Request.Manager', () => {
	let dummyRequestRepo = {
		findById: async function(callback_id: string): Promise<Operations.Request> {
			let r = Deploy.init('michael')
			return Promise.resolve(r)
		},
		store: async function(r: Operations.Request): Promise<Operations.Request> {
			return Promise.resolve(r)
		}
	}
	let dummyNotifier = {
		requestSubmitted(r: Operations.Request) : void {},
		requestApproved(r: Operations.Request) : void {},
		requestDenied(r: Operations.Request) : void {}
	}
	let requestManager = new Request.Manager(dummyRequestRepo, dummyNotifier)

	describe('.listActions()', () => {
		it('returns the expected list', () => {
			let list = requestManager.listActions()
			expect(list).to.have.members(['deploy'])
		})
	})

	describe('.init()', () => {
		let requester = 'michael'

		context('when type = "deploy"', () => {
			it('returns a new Deploy.Request', (done) => {
				requestManager.init('deploy', requester)
				.then(request => {
					expect(request).to.be.an.instanceof(Deploy.Request)
					done()
				})
			})
		})
	})

	describe('.load()', () => {
		it('returns the specified request', (done) => {
			requestManager.load('deploy_michael_12345678')
			.then(request => {
				expect(request).to.be.an.instanceof(Deploy.Request)
				done()
			})
		})
	})

	describe('.update()', () => {
		let r = Deploy.init('michael')
		let s = new Deploy.Server('dev', '123.4.56.789')

		context('when updated with a ResponseURL', () => {
			let rUrl = new Operations.ResponseURL('webhook.slack.com')

			it('returns a list of servers', (done) => {
				requestManager.update(r, rUrl)
				.then(resp => {
					expect(resp).to.eql([s])
					done()
				})
				.catch(done)
			})
			it('and updates the original request', () => {
				expect(r.response_url).to.equal('webhook.slack.com')
			})
		})

		context('when updated with a Server', () => {
			it('returns the updated request', (done) => {
				requestManager.update(r, s)
				.then(resp => {
					expect(resp).to.equal(r)
					done()
				})
			})
			it('updates the original request', () => {
				expect(r.server).to.eql(s)
			})
		})
	})
})
