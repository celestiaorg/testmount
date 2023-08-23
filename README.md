### store.Put safety

This repo demonstrates a safety issue with store.Put that might be causing
block corruption.

branch names correspond to celestia node versions.

# TODO
* report issue
* create issue for defer f.close
    - we should not ignore error from f.close: https://www.joeshaw.org/dont-defer-close-on-writable-files/
* create issue for fsync
    - we should fsync before close so that sync errs are not ignored: https://michael.stapelberg.ch/posts/2017-01-28-golang_atomically_writing/
