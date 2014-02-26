Sensu-to-Alerta
===============

Forward Sensu events to Alerta for a consolidated view and improved visualisation.

Transform this ...

![sensu](/docs/images/sensu.png?raw=true)

Into this ...

![alerta](/docs/images/alerta.png?raw=true)


Installation
------------

To install the Sensu Plugin libary and dependencies...

    $ gem install sensu-plugin httparty

Add the alerta handler and config file...

    $ wget -qO /etc/sensu/handlers/alerta.rb https://raw.github.com/alerta/sensu-alerta/master/alerta.rb
    $ wget -qO /etc/sensu/conf.d/alerta.json https://raw.github.com/alerta/sensu-alerta/master/alerta.json

Configuration
-------------

Replace the example sensu.json.example file...

    $ sudo mv /etc/riemann/riemann.config{,.bak}
    $ wget https://raw.github.com/alerta/riemann-alerta/master/riemann.config
    $ sudo cp riemann.config /etc/riemann/

Restart Sensu...
    
    $ sudo service sensu-server restart


Testing
-------

Generate some test alerts by touching and removing a file called `/ok`...

    $ touch /ok
    $ rm /ok

Vagrant
-------

Alternatively, make use of the [vagrant-try-alerta](https://github.com/alerta/vagrant-try-alerta) repo...

    $ git clone https://github.com/alerta/vagrant-try-alerta.git
    $ cd vagrant-try-alerta
    $ vagrant up alerta-sensu
    $ vagrant ssh alerta-sensu

License
-------

Copyright (c) 2013 Nick Satterly. Available under the MIT License.
