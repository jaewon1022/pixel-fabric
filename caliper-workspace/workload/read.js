"use strict";

const { WorkloadModuleBase } = require("@hyperledger/caliper-core");

class ReadWorkload extends WorkloadModuleBase {
  constructor() {
    super();
  }

  async initializeWorkloadModule(
    workerIndex,
    totalWorkers,
    roundIndex,
    roundArguments,
    sutAdapter,
    sutContext
  ) {
    await super.initializeWorkloadModule(
      workerIndex,
      totalWorkers,
      roundIndex,
      roundArguments,
      sutAdapter,
      sutContext
    );

    for (let i = 1; i <= this.roundArguments.assets; i++) {
      const assetId = `${this.workerIndex}-${i}`;
      const tokenSymbol = `TTN${i}`;

      console.log(
        `Creating asset "${this.workerIndex}-${i}" by workerNode ${workerIndex}`
      );

      let txArgs = {
        contractId: this.roundArguments.contractId,
        contractFunction: "mint",
        invokerIdentity: "client",
        contractArguments: [assetId, tokenSymbol, "100", "user1"],
        readOnly: false,
      };

      await this.sutAdapter.sendRequests(txArgs);
    }
  }

  async submitTransaction() {
    const randomId = Math.floor(Math.random() * this.roundArguments.assets) + 1;

    let txArgs = {
      contractId: this.roundArguments.contractId,
      contractFunction: "queryUser",
      invokerIdentity: "client",
      contractArguments: ["user1"],
      readOnly: true,
    };

    return this.sutAdapter.sendRequests(txArgs);
  }
/*
  async cleanupWorkloadModule() {
      console.log(
        "Deleting all Tokens while testing READ Transaction"
      );

      let txArgs = {
        contractId: this.roundArguments.contractId,
        contractFunction: "deleteAllTokens",
        invokerIdentity: "client",
        contractArguments: [],
        readOnly: false,
      };

      await this.sutAdapter.sendRequests(txArgs);
    }
*/
}

function createWorkloadModule() {
  return new ReadWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
