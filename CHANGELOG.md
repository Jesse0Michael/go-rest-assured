# [4.3.0](https://github.com/Jesse0Michael/go-rest-assured/compare/v4.2.0...v4.3.0) (2024-02-28)

### Features

- remove go-kit/kit and gorilla/mux in favor of net/http's ServeMux ([30da74d](https://github.com/Jesse0Michael/go-rest-assured/commit/30da74da27cea59b735a47a6083502cd734ad10c))

# [4.2.0](https://github.com/Jesse0Michael/go-rest-assured/compare/v4.1.0...v4.2.0) (2024-02-23)

### Features

- add custom http verb support ([cf62869](https://github.com/Jesse0Michael/go-rest-assured/commit/cf62869d6ef1ba9af66039a541c296fd19f55ade))

# [4.1.0](https://github.com/Jesse0Michael/go-rest-assured/compare/v4.0.3...v4.1.0) (2023-12-05)

### Features

- support passing a slog.logger ([902baf0](https://github.com/Jesse0Michael/go-rest-assured/commit/902baf096722ac65a444b6493cd49e9e2630e382))

## [4.0.3](https://github.com/Jesse0Michael/go-rest-assured/compare/v4.0.2...v4.0.3) (2023-11-29)

### Bug Fixes

- skip empty options ([5afb731](https://github.com/Jesse0Michael/go-rest-assured/commit/5afb7319137610ce04d441ea0a0b26cf5d300d2a))

## [4.0.2](https://github.com/Jesse0Michael/go-rest-assured/compare/v4.0.1...v4.0.2) (2023-09-15)

### Chores

- change default branch to main ([3eca768](https://github.com/Jesse0Michael/go-rest-assured/commit/3eca7685cab2b410ecedac6e40c9b63a82e1c30b))

### Continuous Integration

- release from main ([8e01820](https://github.com/Jesse0Michael/go-rest-assured/commit/8e018207b427df9a7eb6b2fc754e611fdd38929a))

## [4.0.1](https://github.com/Jesse0Michael/go-rest-assured/compare/v4.0.0...v4.0.1) (2023-09-15)

### Build System

- move docker labels ([be7180e](https://github.com/Jesse0Michael/go-rest-assured/commit/be7180eeffa882731ade7397e8660d7c738de932))

# [3.2.0](https://github.com/Jesse0Michael/go-rest-assured/compare/v3.1.2...v3.2.0) (2023-09-15)

### Bug Fixes

- use google/uuid package ([f440642](https://github.com/Jesse0Michael/go-rest-assured/commit/f44064258ee2405495f717e34f08892a6d231dc7))

### Build System

- add docker labels ([5abe583](https://github.com/Jesse0Michael/go-rest-assured/commit/5abe58315dd85d4b90d4a7940ce1ba472f399219))

### Chores

- update license ([ebce1b2](https://github.com/Jesse0Michael/go-rest-assured/commit/ebce1b2bac5fe337a8604a563096912527490ebf))
- upgrade to go 1.21 ([7adaa3f](https://github.com/Jesse0Michael/go-rest-assured/commit/7adaa3fbe3ad4ad6db034be15d423fb7776e0b46))

### Continuous Integration

- update go version ([2ea667c](https://github.com/Jesse0Michael/go-rest-assured/commit/2ea667c69a388c85f5bee6164c76d3b39c9d3ac7))

### Features

- add NewClientServe function ([453d258](https://github.com/Jesse0Michael/go-rest-assured/commit/453d25892330d7520f57dd12e3c550cc4ff2285a))
- BREAKING CHANGE upgrade rest assured to v4 ([04f971b](https://github.com/Jesse0Michael/go-rest-assured/commit/04f971bd0d5e96c971bd84a4e7f3e66ecf5e7b62))
- export Serve method to start http listener ([e3af4cc](https://github.com/Jesse0Michael/go-rest-assured/commit/e3af4ccd4e2a911afe6951365f2196ee41c02266))
- move to log/slog for logging ([ac6be67](https://github.com/Jesse0Michael/go-rest-assured/commit/ac6be679ee426bb2fd1235ecb1f401ab2a6b8bd4))

### Tests

- Serve rest assured client in tests ([5503a7d](https://github.com/Jesse0Michael/go-rest-assured/commit/5503a7dc008f9385767c47556473c90dfe922921))

## [3.1.2](https://github.com/Jesse0Michael/go-rest-assured/compare/v3.1.1...v3.1.2) (2022-10-10)

### Continuous Integration

- publish with image name ([8e22407](https://github.com/Jesse0Michael/go-rest-assured/commit/8e22407fb93d6ad3a7742bb2e03a7295cc600e4b))

## [3.1.1](https://github.com/Jesse0Michael/go-rest-assured/compare/v3.1.0...v3.1.1) (2022-10-10)

### Continuous Integration

- update docker publish ([6d8f225](https://github.com/Jesse0Michael/go-rest-assured/commit/6d8f2259b1c9778aeb742f932c81fec781790180))

# [3.1.0](https://github.com/Jesse0Michael/go-rest-assured/compare/v3.0.2...v3.1.0) (2022-10-05)

### Documentation

- update documentation ([1b04b01](https://github.com/Jesse0Michael/go-rest-assured/commit/1b04b01208f670195efce904b2528414eda162cc))

### Features

- add Assured-Method header option ([84e7754](https://github.com/Jesse0Michael/go-rest-assured/commit/84e7754716b11e05d91ef3eccd4d5be0eeb61eab))

## [3.0.2](https://github.com/Jesse0Michael/go-rest-assured/compare/v3.0.1...v3.0.2) (2022-10-05)

### Build System

- build image from static image ([4c29aa3](https://github.com/Jesse0Michael/go-rest-assured/commit/4c29aa35398873ec13afa400650aa6d7499859cd))
- update to go 1.19 ([2ed0296](https://github.com/Jesse0Michael/go-rest-assured/commit/2ed0296a95db1506af7fa78a476727368e960a0c))

### Chores

- update with linter suggestions ([7401f6f](https://github.com/Jesse0Michael/go-rest-assured/commit/7401f6f8b5f94a1bdb237ded38d004e6a8cef9ef))

### Continuous Integration

- update ci actions ([47efa38](https://github.com/Jesse0Michael/go-rest-assured/commit/47efa3886c007830def26aec4887c7b543b22292))

### Tests

- add time for async test ([30dea25](https://github.com/Jesse0Michael/go-rest-assured/commit/30dea2540a8668ad5dcc7078efd1615145bcae71))

## [3.0.1](https://github.com/Jesse0Michael/go-rest-assured/compare/v3.0.0...v3.0.1) (2021-07-30)

### Bug Fixes

- update go mod for v3 ([467b9e5](https://github.com/Jesse0Michael/go-rest-assured/commit/467b9e539cedefb0a0563e3a55a2697c44865367))

# [3.0.0](https://github.com/Jesse0Michael/go-rest-assured/compare/v2.0.11...v3.0.0) (2020-09-03)

### Chores

- update dependencies ([08ad385](https://github.com/Jesse0Michael/go-rest-assured/commit/08ad385b1017e432db395e0f6501b8d62d74e223))
- update for v3 ([d298cc7](https://github.com/Jesse0Michael/go-rest-assured/commit/d298cc76b2c0966fa7eb3572eaa0f5f9f60eefc4))

### Features

- add host option ([6490699](https://github.com/Jesse0Michael/go-rest-assured/commit/649069963e5aeaef88ef649d3eb5eecf43e3a5ce))
- add tls option ([1213d71](https://github.com/Jesse0Michael/go-rest-assured/commit/1213d718729e63536a59d886745c09d173293999))
- remove logger arg ([82f431c](https://github.com/Jesse0Michael/go-rest-assured/commit/82f431cd2b9a95a63b7f08c2cb3d0e5df83d1bc5))
- use functional options ([dedd374](https://github.com/Jesse0Michael/go-rest-assured/commit/dedd374b9b076a5194f70a4fabb2652614804c9f))

### Tests

- fix tls test ([93236a1](https://github.com/Jesse0Michael/go-rest-assured/commit/93236a163aead07bcb80c3e2f00f63f127591b91))

### BREAKING CHANGE

- remove logger option, if you want to redirect your logs in a file, use the appropriate cli commands
- Use funcional options with sane defaults for configuring the assured client instead of a settings struct

## [2.0.11](https://github.com/Jesse0Michael/go-rest-assured/compare/v2.0.10...v2.0.11) (2020-09-01)

### Continuous Integration

- add linter ([efe26b8](https://github.com/Jesse0Michael/go-rest-assured/commit/efe26b8078affaaea85f2550aded81cadfc5df97))

## [2.0.10](https://github.com/Jesse0Michael/go-rest-assured/compare/v2.0.9...v2.0.10) (2020-08-28)

### Documentation

- add response examle ([f0b00d0](https://github.com/Jesse0Michael/go-rest-assured/commit/f0b00d0357e8103efd86c7608c4efbf46198c4c4))

## [2.0.9](https://github.com/Jesse0Michael/go-rest-assured/compare/v2.0.8...v2.0.9) (2020-08-26)

### Continuous Integration

- push latest ([590e4c8](https://github.com/Jesse0Michael/go-rest-assured/commit/590e4c86b9620fe8e9b4ba2fe4e26b4f70da23dc))

## [2.0.8](https://github.com/Jesse0Michael/go-rest-assured/compare/v2.0.7...v2.0.8) (2020-08-26)

### Continuous Integration

- set latest tag ([104f1eb](https://github.com/Jesse0Michael/go-rest-assured/commit/104f1ebac7c0f4149cbc922d375b075f7cb8bcae))

## [2.0.7](https://github.com/Jesse0Michael/go-rest-assured/compare/v2.0.6...v2.0.7) (2020-08-26)

### Continuous Integration

- separate tagging and releasing ([d0464f6](https://github.com/Jesse0Michael/go-rest-assured/commit/d0464f6cb0c518fd15193f32138c66c11d5b48d0))

## [2.0.6](https://github.com/Jesse0Michael/go-rest-assured/compare/v2.0.5...v2.0.6) (2020-08-26)

### Continuous Integration

- include release reference ([3e7eae4](https://github.com/Jesse0Michael/go-rest-assured/commit/3e7eae410bc2ecfd567eb795d5c4a7433a113dcd))

### Other

- Merge branch 'master' of ssh://github.com/Jesse0Michael/go-rest-assured ([137c9fd](https://github.com/Jesse0Michael/go-rest-assured/commit/137c9fdcb69531310a5cb47d8c77e629f26bc424))

## [2.0.5](https://github.com/Jesse0Michael/go-rest-assured/compare/v2.0.4...v2.0.5) (2020-08-26)

### Continuous Integration

- docker build/push ([c18d7a9](https://github.com/Jesse0Michael/go-rest-assured/commit/c18d7a9d07c9682e516a5359fce8c00a6367480e))

## [2.0.4](https://github.com/Jesse0Michael/go-rest-assured/compare/v2.0.3...v2.0.4) (2020-04-15)

### Documentation

- document application preload specification ([5cd071b](https://github.com/Jesse0Michael/go-rest-assured/commit/5cd071b05d9eeb39423f878811cb7c6ed6f56423))

## [2.0.3](https://github.com/Jesse0Michael/go-rest-assured/compare/v2.0.2...v2.0.3) (2020-04-11)

### Documentation

- update build badge ([8a5872b](https://github.com/Jesse0Michael/go-rest-assured/commit/8a5872b192e9e7991e116301debb454ba1c3a627))

## [2.0.2](https://github.com/Jesse0Michael/go-rest-assured/compare/v2.0.1...v2.0.2) (2020-04-11)

### Bug Fixes

- remove old assured application ([355b32e](https://github.com/Jesse0Michael/go-rest-assured/commit/355b32ef74386850afc13d00297e304226e41002))

### Other

- Merge branch 'master' of ssh://github.com/Jesse0Michael/go-rest-assured ([8a82f0b](https://github.com/Jesse0Michael/go-rest-assured/commit/8a82f0bcc45657511d3bef0de0d14a78f677d5b5))

## [2.0.1](https://github.com/Jesse0Michael/go-rest-assured/compare/v2.0.0...v2.0.1) (2020-04-11)

### Bug Fixes

- separeate application and client readme ([b856342](https://github.com/Jesse0Michael/go-rest-assured/commit/b856342b3e8deac048c33cf8381447e01db33628))

# [1.1.0](https://github.com/Jesse0Michael/go-rest-assured/compare/v1.0.1...v1.1.0) (2020-04-10)

### Bug Fixes

- update test package ([122d28f](https://github.com/Jesse0Michael/go-rest-assured/commit/122d28fd1ba47b808e220506479bcbbe6710ad53))

### Features

- reorgonize pkg and cmd directories ([2c00b9a](https://github.com/Jesse0Michael/go-rest-assured/commit/2c00b9af16367919513556a2328b0cacebb396ac))

## [1.0.1](https://github.com/Jesse0Michael/go-rest-assured/compare/v1.0.0...v1.0.1) (2020-04-10)

### Chores

- remove releaserc ([8d2b70b](https://github.com/Jesse0Michael/go-rest-assured/commit/8d2b70b9947f74a0be915ca272d528a7d563776b))

### Other

- Merge branch 'master' of ssh://github.com/Jesse0Michael/go-rest-assured ([804f6f1](https://github.com/Jesse0Michael/go-rest-assured/commit/804f6f1329407eebea46d36d36b6ccee89ae935a))

# 1.0.0 (2020-04-10)

### Bug Fixes

- add github build and test actions ([520e814](https://github.com/Jesse0Michael/go-rest-assured/commit/520e8140ecb11683f7bc2416a107f36ece1e026e))
- fix semantic release ([80ef437](https://github.com/Jesse0Michael/go-rest-assured/commit/80ef437e8aacc0512a19ca80669e918daffcb1e2))
- try codfish's action ([08b1e86](https://github.com/Jesse0Michael/go-rest-assured/commit/08b1e8645a6867f9b02ef5fcc5f73e1428afc46c))
- try release me action ([278bd8f](https://github.com/Jesse0Michael/go-rest-assured/commit/278bd8f0a9c185c189a34d143c474672f02c5ecd))

### Other

- Fix readme typo ([2829caa](https://github.com/Jesse0Michael/go-rest-assured/commit/2829caa783abe1469cabc30f44847495586245d7))
- Fix client test. ([550c064](https://github.com/Jesse0Michael/go-rest-assured/commit/550c06415410d89c8f0f92a3be4d158eab52524f))
- Go module. ([001aabf](https://github.com/Jesse0Michael/go-rest-assured/commit/001aabf84d4a1884f1c61bf934e02de0ba6455b6))
- circle and coverage ([c3d913a](https://github.com/Jesse0Michael/go-rest-assured/commit/c3d913a29414ab1851a52589f0fc99875463d66a))
- Circle CI, my guy (#7) ([57d7b08](https://github.com/Jesse0Michael/go-rest-assured/commit/57d7b085e82ec7911437f96c2ef9e449bb26db36)), closes [#7](https://github.com/Jesse0Michael/go-rest-assured/issues/7)
- write header after having configured it ([d3aa03f](https://github.com/Jesse0Michael/go-rest-assured/commit/d3aa03fea009566c0c0b416adc60a386705f2c5f))
- write header after having configured it ([c56c003](https://github.com/Jesse0Michael/go-rest-assured/commit/c56c0034e660aaf4a382438a41ac07515764f412))
- delay response ([d34f5ab](https://github.com/Jesse0Michael/go-rest-assured/commit/d34f5ab9ac8da5601c6e34878549a7a47bef7957))
- calls track query ([514778e](https://github.com/Jesse0Michael/go-rest-assured/commit/514778e6b6248b251e4d59173c2ccc4dd655f01c))
- Query in call ([b98fe2e](https://github.com/Jesse0Michael/go-rest-assured/commit/b98fe2ebc0f5a524115fbee687e31594de25f0ec))
- Merge branch 'master' into query-in-call ([a169d1b](https://github.com/Jesse0Michael/go-rest-assured/commit/a169d1b144d5d0d16980ebe7c3b24d49d59054b5))
- add query params to call ([08c107e](https://github.com/Jesse0Michael/go-rest-assured/commit/08c107eafee7add16e40700c2337616271758d09))
- custom unmarshalling ([d965d17](https://github.com/Jesse0Michael/go-rest-assured/commit/d965d1731b63803604884ab4d81d23d47c1a2c39))
- custom unmarshalling ([1b3a4b2](https://github.com/Jesse0Michael/go-rest-assured/commit/1b3a4b2603bd5715b7cd54bb56026ec309d9c400))
- tracking calls should be on be default ([841b442](https://github.com/Jesse0Michael/go-rest-assured/commit/841b4425977c6480d1ba989ed5f3e22f18af2028))
- Callbacks ([51dd7ee](https://github.com/Jesse0Michael/go-rest-assured/commit/51dd7ee9d182b231a72b1c9231f6dac81be1dc66))
- callbacks ([31f99bc](https://github.com/Jesse0Michael/go-rest-assured/commit/31f99bc5ff8104077799f244ebab36f3481818e1))
- delay from client ([d14b1fa](https://github.com/Jesse0Michael/go-rest-assured/commit/d14b1fa08e94c2fcd3b2ea4e5000811e0655d083))
- callbacks and tests ([41cea71](https://github.com/Jesse0Michael/go-rest-assured/commit/41cea716dc9acfb0f38c20df02d12e5077c3ee8d))
- callbacks ([b215ae6](https://github.com/Jesse0Michael/go-rest-assured/commit/b215ae6d2cdc81658ecfd3b8e316918f29b89ec8))
- rearrange headers ([d7cbed2](https://github.com/Jesse0Michael/go-rest-assured/commit/d7cbed29d3f6aa598fbed6509e2d78cce69325d1))
- sanitize path on client ([fec602d](https://github.com/Jesse0Michael/go-rest-assured/commit/fec602d610a4a06b8f60b2ff19c7c0c311a20652))
- support headers ([b23211b](https://github.com/Jesse0Michael/go-rest-assured/commit/b23211bdb7881623e491bd10076118f13d2f90cb))
- make tracking made calls disable able ([c8af462](https://github.com/Jesse0Michael/go-rest-assured/commit/c8af462151c0d53bc512d7537209657c49f4481e))
- preload calls ([0824afb](https://github.com/Jesse0Michael/go-rest-assured/commit/0824afbb509f19cd4895fe88fb83e74db3caf61f))
- add port and logger support for binary ([b02c3a9](https://github.com/Jesse0Michael/go-rest-assured/commit/b02c3a98bcfd4b776543b0e45e119ffbbdd64764))
- Allow mocking calls at root level ([2a38c91](https://github.com/Jesse0Michael/go-rest-assured/commit/2a38c91a33e92bc2264dd933966da80a870f4614))
- Allow mocking calls at root level ([f904cf0](https://github.com/Jesse0Michael/go-rest-assured/commit/f904cf0ca1a4d17c3ae597dc1cf6f44b4f471959))
- remove freeport ([e115920](https://github.com/Jesse0Michael/go-rest-assured/commit/e1159205631e2206d76c20d3ae4e266b0eca63ca))
- use sync.Mutex instead of channels ([c24e9cc](https://github.com/Jesse0Michael/go-rest-assured/commit/c24e9ccea35cf7f59b8b0fb8381ad7b9ba939911))
- rename files ([1c85984](https://github.com/Jesse0Michael/go-rest-assured/commit/1c85984a69f54bd32f879c429e826e6f26413546))
- why use pointers here ([3a77f19](https://github.com/Jesse0Michael/go-rest-assured/commit/3a77f19975234a155b1c78ce1181c9207b849251))
- close client ([39309a7](https://github.com/Jesse0Michael/go-rest-assured/commit/39309a785c5e74ed0e996b0518156444d231a166))
- lock assured interactions ([ddb2a6d](https://github.com/Jesse0Michael/go-rest-assured/commit/ddb2a6d18a66a65f88eb2f24a0dba18ba18c12d4))
- Configure or Default Client with Serve ([622b195](https://github.com/Jesse0Michael/go-rest-assured/commit/622b19556bb3459a9ab9b91279a9f0aba57f29c4))
- client usage testing ([81e0925](https://github.com/Jesse0Michael/go-rest-assured/commit/81e0925a1f436f7f5e12e5373bf9779474e3d8a0))
- assured client testing ([0bcebcb](https://github.com/Jesse0Michael/go-rest-assured/commit/0bcebcb8c4a7c2caca2712e46350aa95b7071555))
- client up ([fd4e504](https://github.com/Jesse0Michael/go-rest-assured/commit/fd4e504a85c8acea6738fce14bdb4ebfdaaee6f2))
- change then to verify, and use free port ([4fe7037](https://github.com/Jesse0Michael/go-rest-assured/commit/4fe70370520dd498577e35efe5a614837ded2e32))
- license ([99921e8](https://github.com/Jesse0Michael/go-rest-assured/commit/99921e8453f517262a10fdb0da7e299a61df5187))
- go fmt ([0bb6436](https://github.com/Jesse0Michael/go-rest-assured/commit/0bb6436237ec2c907902220db35da45cba6cbccc))
- move to one package ([15b8587](https://github.com/Jesse0Michael/go-rest-assured/commit/15b85875ff1077eba96eb8af5bd2ed8675348b95))
- populate readme ([9c2d0ce](https://github.com/Jesse0Michael/go-rest-assured/commit/9c2d0ce64db27f8dd70fa578374824e09bbb178d))
- fix binding and endpoints ([7e3bf63](https://github.com/Jesse0Michael/go-rest-assured/commit/7e3bf63df07b1557e1c3cd82aadfc88c92cc6b7a))
- wrap endpoints, test bindings, assure calls ([1e20981](https://github.com/Jesse0Michael/go-rest-assured/commit/1e20981fb494b6134020bc63694de22717cf47b0))
- rest assured endpoints with tests ([02786fa](https://github.com/Jesse0Michael/go-rest-assured/commit/02786fae69287d2d502b821dd495306486fe7b0a))
- assured bindings and endpoints ([9b13b09](https://github.com/Jesse0Michael/go-rest-assured/commit/9b13b091ead860dccec7da61693ce32b1259323e))
- go-rest-assured ([2022320](https://github.com/Jesse0Michael/go-rest-assured/commit/20223202b5097da4536544fc1ee959531dc6b152))
