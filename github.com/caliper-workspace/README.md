#### Caliper Workspace
Pixelterior 프로젝트의 fabric 체인코드 성능 측정을 위해 만들어진 디렉토리입니다.

#### Work Tree
```js
caliper-workspace/
├── benchmarks
│   └── myAssetBenchmark.yaml
├── generated
│   └── // 네이버 클라우드 플랫폼 key files
├── networks
│   └── networkConfig.yaml
├── package.json
├── package-lock.json
├── report.html
└── workload
    ├── read.js
    └── write.js
```
#### 도중에 키 파일이 변경된 경우
다음의 방법을 이용해 새로운 키 파일을 인증 수단으로 이용할 수 있습니다.
1. 네이버 클라우드 플랫폼의 BlockchainService/Nodes 에 접근
2. `Orderers` | `Peers` | `CAs` 탭 중 Cas 탭을 선택
3. `org-ca` 를 선택한 후 사용자 ID 관리 버튼 클릭
   <br/> 3-1. 사용자 ID 추가 ( 추가 시 사용자 유형을 client로 변경해야 함 )
        ![naver-cloud-platform-add-new-keyfile](https://github.com/user-attachments/assets/ce3ae70b-0a5c-4680-bfee-cd904974791b)
   <br/> 3-2. 사용자 ID 인증서 다운로드 ( `인증서 + Key (JSON)` 버튼 클릭 )
        ![naver-cloud-platform-download-info](https://github.com/user-attachments/assets/95b2e329-c82d-41e0-b3f4-547e3fc5d489)
4. 다운로드된 인증서를 `generated/` 폴더에 추가
5. 다운로드된 파일을 base64 형식으로 변환
```
jq -r .key {파일 위치} | base64 -d > generated/user.key
jq -r .cert {파일 위치} | base64 -d > generated/user.cert
```
6. 생성된 키파일 사용
