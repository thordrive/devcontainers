# TODO

## Collision Cehck

```yml
name: thordrive/dev-gcc
images:
  - tags: ["11.2.0", "11.2", "11"]
    from: "library/gcc:11.2.0"

  - tags: ["11.1.0", "11.1", "11"]
    from: "library/gcc:11.1.0"
```

Reference `thordrive/dev-gcc:11` is collide from `library/gcc:11.2.0` and `library/gcc.11.1.0`. It must throw an error.


## Pattern Matching `partially done`

```yml
images:
  - tags: ["16.{0}.{1}", "16.{0}", "16"]
    from: "library/node:/16\.(\d+)\.(\d+)-bullseye/"
    # strategy: "semver"
```
If the tag starts and ends with `/`, the tag will be evaluated as regex and pattern matching is enabled.

Above manifest will match tags `16.13.0`, `16.13.1`, `16.14.0`... from DockerHub and generate tags `16.13.0`, `16.13.1`, `16.14`, `16.13`, `16.14`, and `16`.
  
Only *semver* is supported on current implementation. *semver* is a default value of `strategy` field which is does not exists on current implementation.

In the semver strategy, the larger the number captured earlier, the higher the priority. Duplicate tags with lower priority are deleted.


### Issue

- [x] How to choose base image of `16.13` or `16`?
- [ ] Strategy to resolve codename (e.g. ROS)
- [ ] Improve *semver* strategy to bypass the captured strings.
