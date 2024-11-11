"use strict";

const { WorkloadModuleBase } = require("@hyperledger/caliper-core");

class ReadWorkload extends WorkloadModuleBase {
  constructor() {
    super();
  }

  async submitTransaction() {
    let txArgs = {
      contractId: this.roundArguments.contractId,
      contractFunction: "queryUser",
      invokerIdentity: "client",
      contractArguments: ["user1"],
      readOnly: true,
    };

    return this.sutAdapter.sendRequests(txArgs);
  }
}

function createWorkloadModule() {
  return new ReadWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
