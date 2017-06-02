import * as Operations from '../operations/operations'
import * as Deploy from '../operations/deploy'

export class Manager {
	constructor(public requestRepo: Operations.RequestRepo) {}

	async listActions() : Promise<Array<string>> {
		return this.requestRepo.findAllTypes()
	}

	async init(type: string, requester: string) : Promise<Operations.Request> {
		let r = Deploy.init(requester)
		return this.requestRepo.store(r)
	}

	async load(id: number) : Promise<Operations.Request> {
		return this.requestRepo.findById(id)
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