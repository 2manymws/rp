# Changelog

## [v0.8.3](https://github.com/2manymws/rp/compare/v0.8.2...v0.8.3) - 2023-12-15

## [v0.8.2](https://github.com/2manymws/rp/compare/v0.8.1...v0.8.2) - 2023-12-15

## [v0.8.1](https://github.com/k1LoW/rp/compare/v0.8.0...v0.8.1) - 2023-12-15
### Fix bug 🐛
- Don't touch req.URL by @k1LoW in https://github.com/k1LoW/rp/pull/38

## [v0.8.0](https://github.com/k1LoW/rp/compare/v0.7.0...v0.8.0) - 2023-11-30
### New Features 🎉
- support error handler by @pyama86 in https://github.com/k1LoW/rp/pull/35

## [v0.7.0](https://github.com/k1LoW/rp/compare/v0.6.0...v0.7.0) - 2023-11-24
### New Features 🎉
- Make it possible to proxy through a Unix domain socket. by @pyama86 in https://github.com/k1LoW/rp/pull/33

## [v0.6.0](https://github.com/k1LoW/rp/compare/v0.5.0...v0.6.0) - 2023-11-01
### New Features 🎉
- Add a hook for errors occurring before the roundtrip. by @pyama86 in https://github.com/k1LoW/rp/pull/32
### Other Changes
- Set the number of CPUs to worker_processs by @k1LoW in https://github.com/k1LoW/rp/pull/30

## [v0.5.0](https://github.com/k1LoW/rp/compare/v0.4.5...v0.5.0) - 2023-10-12
### Breaking Changes 🛠
- Minimize the number of methods that must be implemented. by @k1LoW in https://github.com/k1LoW/rp/pull/28
- To act as a minimum reverse proxy with only the GetUpstream implementation, make it execute SetXForwarded(). by @k1LoW in https://github.com/k1LoW/rp/pull/29
### Other Changes
- Run 2 benchmarks on same runner by @k1LoW in https://github.com/k1LoW/rp/pull/26

## [v0.4.5](https://github.com/k1LoW/rp/compare/v0.4.4...v0.4.5) - 2023-09-29

## [v0.4.4](https://github.com/k1LoW/rp/compare/v0.4.3...v0.4.4) - 2023-09-29
### Fix bug 🐛
- Should change pr.Out.Host by @k1LoW in https://github.com/k1LoW/rp/pull/24
### Other Changes
- Add gostyle-action by @k1LoW in https://github.com/k1LoW/rp/pull/23

## [v0.4.3](https://github.com/k1LoW/rp/compare/v0.4.2...v0.4.3) - 2023-09-04
### Fix bug 🐛
- Fix cache setting of NGINX by @k1LoW in https://github.com/k1LoW/rp/pull/21

## [v0.4.2](https://github.com/k1LoW/rp/compare/v0.4.1...v0.4.2) - 2023-09-04
### New Features 🎉
- Enable ngx_http_js_module in upstream server by @k1LoW in https://github.com/k1LoW/rp/pull/19
### Other Changes
- Test using HTTP/2 by @k1LoW in https://github.com/k1LoW/rp/pull/16
- Show benchmark in pull request using octocov by @k1LoW in https://github.com/k1LoW/rp/pull/17
- Freeze benchtime by @k1LoW in https://github.com/k1LoW/rp/pull/18

## [v0.4.1](https://github.com/k1LoW/rp/compare/v0.4.0...v0.4.1) - 2023-08-27

## [v0.4.0](https://github.com/k1LoW/rp/compare/v0.3.0...v0.4.0) - 2023-08-26
### Other Changes
- Clean up go.mod by @k1LoW in https://github.com/k1LoW/rp/pull/11

## [v0.3.0](https://github.com/k1LoW/rp/compare/v0.2.0...v0.3.0) - 2023-08-26
### New Features 🎉
- Add ListenAndServe and ListenAndServeTLS for convenience. by @k1LoW in https://github.com/k1LoW/rp/pull/10
### Other Changes
- Set up an environment for enhanced benchmarking by @k1LoW in https://github.com/k1LoW/rp/pull/5
- Bump github.com/docker/docker from 20.10.7+incompatible to 20.10.24+incompatible by @dependabot in https://github.com/k1LoW/rp/pull/6
- Use b.RunParallel to run benchmarks in parallel. by @k1LoW in https://github.com/k1LoW/rp/pull/8
- Use NewServer and NewTLSServer instead of httptest.* by @k1LoW in https://github.com/k1LoW/rp/pull/9

## [v0.2.0](https://github.com/k1LoW/rp/compare/v0.1.0...v0.2.0) - 2023-08-25
### Breaking Changes 🛠
- Do not add X-Forwarded-* header by default by @k1LoW in https://github.com/k1LoW/rp/pull/4
### Other Changes
- add permance test with nginx by @pyama86 in https://github.com/k1LoW/rp/pull/3

## [v0.1.0](https://github.com/k1LoW/rp/commits/v0.1.0) - 2023-08-25
