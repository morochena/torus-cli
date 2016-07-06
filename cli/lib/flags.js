/**
 * Utility library for standardizing flags, their naming, and the defaults.
 */

'use strict';

var Command = require('./cli/command');

var flags = exports;

var LIST = flags.LIST = {
  force: {
    usage: '-f, --force',
    description: 'force override',
    default: undefined
  },
  project: {
    usage: '-p, --project [name]',
    description: 'specify the project',
    default: undefined
  },
  org: {
    usage: '-o, --org [name]',
    description: 'specify the organization',
    default: undefined
  },
  service: {
    usage: '-s, --service [name]',
    description: 'specify the service',
    default: undefined
  },
  environment: {
    usage: '-e, --environment [name]',
    description: 'specify the environment',
    default: undefined
  },
  instance: {
    usage: '-i, --instance [name]',
    description: 'specify the instance',
    default: '1'
  }
};

flags.add = function (cmd, name, overrides) {
  if (!(cmd instanceof Command)) {
    throw new Error('Must provide an instance of a Command');
  }

  if (!LIST[name]) {
    throw new Error('Unknown option: ' + name);
  }

  var matching = cmd.options.filter(function (o) {
    return (o.name() === name);
  });

  if (matching.length > 0) {
    throw new Error('Cannot add the same option twice');
  }

  overrides = overrides || {};

  var opt = LIST[name];
  cmd.option.apply(cmd, [
    opt.usage,
    (overrides.description) ? overrides.description : opt.description,
    (overrides.default) ? overrides.default : opt.default
  ]);
};
