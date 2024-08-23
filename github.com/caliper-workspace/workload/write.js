'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class MyWorkload extends WorkloadModuleBase {
    constructor() {
        super();
    }

    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);

        for (let i=0; i<this.roundArguments.assets; i++) {
            const assetID = `${this.workerIndex}_${i}`;
            console.log(`Worker ${this.workerIndex}: Creating asset ${assetID}`);
            const request = {
                contractId: this.roundArguments.contractId,
                contractFunction: 'createAsset',
                invokerIdentity: 'client',
                contractArguments: [assetID,'100', '10000'],
                readOnly: false
            };

            await this.sutAdapter.sendRequests(request);
        }
    }

    async submitTransaction() {
        const randomId = Math.floor(Math.random()*this.roundArguments.assets);

	const assetSet = [
	    { price: "10", totalStock: "100000" },
	    { price: "20", totalStock: "50000" },
	    { price: "40", totalStock: "25000" },
	    { price: "50", totalStock: "20000" },
	    { price: "100", totalStock: "10000" },
	    { price: "200", totalStock: "5000" },
	];

	const randomAsset = assetSet[Math.floor(Math.random() * assetSet.length)];

        const myArgs = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'updateAsset',
            invokerIdentity: 'client',
            contractArguments: [`${this.workerIndex}_${randomId}`, randomAsset.price, randomAsset.totalStock],
            readOnly: false
        };

        await this.sutAdapter.sendRequests(myArgs);
    }

    async cleanupWorkloadModule() {
        for (let i=0; i<this.roundArguments.assets; i++) {
            const assetID = `${this.workerIndex}_${i}`;
            console.log(`Worker ${this.workerIndex}: Deleting asset ${assetID}`);
            const request = {
                contractId: this.roundArguments.contractId,
                contractFunction: 'deleteAsset',
                invokerIdentity: 'client',
                contractArguments: [assetID],
                readOnly: false
            };

            await this.sutAdapter.sendRequests(request);
        }
    }
}

function createWorkloadModule() {
    return new MyWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
