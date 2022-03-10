# myhttp

**myhttp** is a tool which makes http requests in parallel and prints the address of the request along with the
MD5 hash of the response.

### Examples 

```
$> myhttp http://www.adjust.com http://google.com

http://www.adjust.com d1b40e2a2ba488a054186e4ed0733f9752f66949
http://google.com 9d8ec921bdd275fb2a605176582e08758eb60641
```

Concurrency limit can be provided through the argument `--parallel`, 
which equals 10 by default:

```
$> myhttp --parallel 2 http://www.adjust.com http://google.com

http://www.adjust.com d1b40e2a2ba488a054186e4ed0733f9752f66949
http://google.com 9d8ec921bdd275fb2a605176582e08758eb60641
```

### Installation 
`go install`