"use strict";

const { WorkloadModuleBase } = require("@hyperledger/caliper-core");

class WriteWorkload extends WorkloadModuleBase {
  constructor() {
    super();
  }

  async submitTransaction() {
    const crypto = require("crypto");
    const id = crypto.randomBytes(16).toString("hex");

    const txArgs = {
      contractId: this.roundArguments.contractId,
      contractFunction: "mint",
      invokerIdentity: "client",
      contractArguments: [id, "100"],
      readOnly: false,
    };

    return this.sutAdapter.sendRequests(txArgs);
  }
}

function createWorkloadModule() {
  return new WriteWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;

