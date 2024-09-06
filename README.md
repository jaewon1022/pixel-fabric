### Pixelterior-Fabric

Pixelterior 프로젝트의 P2P 거래 기능 구현을 위해 `hyperledger-fabric` 과 `go` 를 이용해 만들어진 프로젝트 입니다.

기본적인 서버 구성은 [네이버 클라우드 플랫폼](https://console.ncloud.com/dashboard)을 이용해 설정되었으며, 서버의 스펙은 다음과 같습니다.

#### Server Spec
해당 프로젝트는 네이버 클라우드 플랫폼 `high-cpu c16-g3` 요금제를 사용하고 있으며,

각 서버 스펙 별 자세한 요금제는 [여기](https://www.ncloud.com/product/compute/ssdServer#pricing)에서 확인하실 수 있습니다.

```
Server Name : pixel-nodepool-w-5jeq (26161172)
Server Image : Ubuntu-22.04-nksw
vCPU : 16 개
Memory : 32 GB
Fabric Version : 1.4.12
Caliper Version : 0.6.0
Node Version : v12.22.12
Npm Version : v6.14.16
```

#### 체인 코드 수정 및 배포 방법
1. docker 접속 ( docker run -ti --rm -v $GOPATH/src:/opt/gopath/src hyperledger/fabric-tools:1.4.12 /bin/bash )
2. 체인 코드 수정 ( github.com/hyperledger/fabric/examples/chaincode/go/pixelterior/chaincode.go )
3. /cds 파일로 이동
4. 체인코드 패키징 ( peer chaincode package -n `pixelterior-go` -p `github.com/hyperledger/fabric/examples/chaincode/go/pixelterior/cmd` -v `0.0.000` `pixelterior-go.v0.0.000.cds` )
5. 생성된 cds 파일 커밋 후 네이버 클라우드 플랫폼에 cds 업로드 및 인스턴스화

## Docs
### ERC-1155 표준?
- <a href="https://andyjaewon.notion.site/PixelTerior-ERC-1155-c98676a0a3ae46ce85f1ec86d2815b7c">ERC-1155 표준</a>
### UTXO vs Account Model
- <a href="UTXO%20vs%20Account%20Model.md">UTXO vs Account Model</a>
