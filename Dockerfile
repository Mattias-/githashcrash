FROM scratch
ADD bin/githashcrash /githashcrash
ENTRYPOINT ["/githashcrash"]
