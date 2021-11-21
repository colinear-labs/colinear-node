FROM electriccoinco/zcashd

COPY --chown=zcashd:zcashd ./zcash.conf /srv/zcashd/.zcash/zcash.conf
