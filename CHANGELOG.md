# Change Log
All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/).

## [1.4.1] - 2016-08-26
### Added
- #34: Added "BOUNDS key" command

### Fixed
- #38: Allow for nginx support
- #39: Reset requirepass 

## [1.3.0] - 2016-07-22
### Added
- New EXPIRE, PERSISTS, TTL commands. New EX keyword to SET command.
- Support for plain strings using `SET ... STRING value.` syntax.
- New SEARCH command for finding strings. 
- Scans can now order descending.

### Fixed
- #28: fix windows cli issue

## [1.2.0] - 2016-05-24
### Added
- #17: Roaming Geofences for NEARBY command
- #15: maxmemory config setting

## [1.1.4] - 2016-04-19
### Fixed
- #12: Issue where a newline was being added to HTTP POST requests.
- #13: OBJECT keyword not accepted for WITHIN command
- Panic on missing key for search requests.

## [1.1.2] - 2016-04-12
### Fixed
- A glob suffix wildcard can result in extra hits.
- The native live geofence sometimes fails connections.

## [1.1.0] - 2016-04-02
### Added
- Resp client support. All major programming languages now supported.
- Added WITHFIELDS option to GET.
- Added OUTPUT command to allow for outputing JSON when using RESP.
- Added DETECT option to geofences.

### Changes
- New AOF file structure.
- Quicker and safer AOFSHRINK.

### Deprecation Warning
- Native protocol support is being deprecated in a future release in favor of RESP.
