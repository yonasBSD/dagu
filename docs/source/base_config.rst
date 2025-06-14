.. _base configuration:

Base Configuration for all DAGs
=====================================

Creating a base configuration (default path: ``~/.config/dagu/base.yaml``) is a convenient way to organize shared settings among all DAGs. The path to the base configuration file can be configured. See :ref:`Configuration Options` for more details.

Any settings defined in the base configuration can be overridden by individual DAG files. This is particularly useful for queue management and other organizational defaults.

Example:

.. code-block:: yaml

    # Queue management settings - shared by all DAGs
    queue: "global-queue"      # Assign all DAGs to the same queue
    maxActiveRuns: 2           # Default concurrent runs per DAG

    # directory path to save logs from standard output
    logDir: /path/to/stdout-logs/

    # history retention days (default: 30)
    histRetentionDays: 3

    # Email notification settings
    mailOn:
      failure: true
      success: true

    # SMTP server settings
    smtp:
      host: "smtp.foo.bar"
      port: "587"
      username: "<username>"
      password: "<password>"

    # Error mail configuration
    errorMail:
      from: "foo@bar.com"
      to: "foo@bar.com"
      prefix: "[Error]"

    # Info mail configuration
    infoMail:
      from: "foo@bar.com"
      to: "foo@bar.com"
      prefix: "[Info]"