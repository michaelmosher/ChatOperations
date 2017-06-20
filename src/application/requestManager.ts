import * as Operations from '../operations/operations'
import * as Deploy from '../operations/deploy'

interface notifier {
	requestSubmitted(r: Operations.Request) : void
	requestApproved(r: Operations.Request) : void
	requestDenied(r: Operations.Request) : void
}

export class Manager {
	constructor(public requestRepo: Operations.RequestRepo, public notifier: notifier) {}

	listActions() : Array<string> {
		return ['deploy']
	}

	async init(type: string, requester: string) : Promise<Operations.Request> {
		let r = Deploy.init(requester)
		return this.requestRepo.store(r)
	}

	async load(callback_id: string) : Promise<Operations.Request> {
		return this.requestRepo.findById(callback_id)
	}

	async update(r: Operations.Request, u: any) : Promise<Operations.Request|Array<any>> {
		r.update(u)

		return this.requestRepo.store(r)
		.then(() => {
			if (r.isReady()) {
				return r
			} else {
				return r.next()
			}
		})
	}

	approveRequest(r: Operations.Request) : void {
		this.update(r, new Operations.Approval(true))
		// TODO send notification
	}

	rejectRequest(r: Operations.Request) : void {
		this.update(r, new Operations.Approval(false))
		// TODO send notification
	}

	executeRequest(r: Operations.Request) : void {
		// TODO
	}
}