# Solidity API

## ZkEvmV2

### MODULO_R

```solidity
uint256 MODULO_R
```

### OPERATOR_ROLE

```solidity
bytes32 OPERATOR_ROLE
```

### currentTimestamp

```solidity
uint256 currentTimestamp
```

_DEPRECATED in favor of currentFinalizedState hash._

### currentL2BlockNumber

```solidity
uint256 currentL2BlockNumber
```

The most recent finalized L2 block number.

### stateRootHashes

```solidity
mapping(uint256 => bytes32) stateRootHashes
```

The most recent L2 state root hash mapped by block number.

### verifiers

```solidity
mapping(uint256 => address) verifiers
```

The verifier address to use for a proof type when proving.

### _verifyProof

```solidity
function _verifyProof(uint256 _publicInput, uint256 _proofType, bytes _proof) internal
```

Verifies the proof with locally computed public inputs.

_If the verifier based on proof type is not found, it reverts with InvalidProofType._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _publicInput | uint256 | The computed public input hash cast as uint256. |
| _proofType | uint256 | The proof type to determine which verifier contract to use. |
| _proof | bytes | The proof to be verified with the proof type verifier contract. |
