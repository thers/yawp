// @flow
import test from 'test';

export const fn = async (...a) => {
    const [as, ...r] = a;

    return as * test;
};
