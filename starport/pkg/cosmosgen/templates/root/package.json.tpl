{
  "name": "{{ .User }}-{{ .Repo }}-ts-client",
  "version": "0.0.1",
  "description": "Autogenerated Typescript Client",
  "author": "Ignite Codegen <hello@ignt.com>",
  "license": "Apache-2.0",
  "licenses": [
    {
      "type": "Apache-2.0",
      "url": "http://www.apache.org/licenses/LICENSE-2.0"
    }
  ],
  "main": "index.ts",
  "publishConfig": {
    "access": "public"
  },
  "dependencies": {
    "@cosmjs/launchpad": "0.27.0",
    "@cosmjs/proto-signing": "0.27.0",
    "@cosmjs/stargate": "0.27.0",
    "axios": "^0.25.0",
    "buffer": "^6.0.3",
    "reconnecting-websocket": "^4.4.0",
    "eventemitter3": "4.0.7"
  },
  "devDependencies": {
    "@keplr-wallet/types": "^0.9.10"
  },
  "peerDependencies": {
    "@cosmjs/launchpad": "0.27.0",
    "@cosmjs/proto-signing": "0.27.0",
    "@cosmjs/stargate": "0.27.0"
  }
}