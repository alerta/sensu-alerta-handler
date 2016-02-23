#!/usr/bin/env ruby

require 'rubygems' if RUBY_VERSION < '1.9.0'
require 'sensu-handler'
require 'socket'
require 'httparty'
require 'json'
require 'timeout'


class Alerta < Sensu::Handler

  def short_name
    @event['client']['name'] + '/' + @event['check']['name']
  end

  def action_to_string
    @event['action'].eql?('resolve') ? "RESOLVED" : "ALERT"
  end

  def status_to_severity
    case @event['check']['status']
      when 0
        "normal"
      when 1
        "warning"
      when 2
        "critical"
      else
        "unknown"
    end
  end

  def handle
    endpoint = settings['alerta']['endpoint'] || 'http://localhost:8080'
    key = settings['alerta']['key'] || nil

    url = endpoint + '/alert'
    puts url
    hostname = Socket.gethostname

    environment = @event['check']['environment'] || 'Production'
    subscribers = @event['check']['subscribers'] || []

    payload = {
      "origin" => "sensu/#{hostname}",
      "resource" => "#{@event['client']['name']}:#{@event['client']['address']}",
      "event" => "#{@event['check']['name']}",
      "group" => "Sensu",
      "severity" => "#{status_to_severity}",
      "environment" => environment,
      "service" => @event['client']['subscriptions'],
      "tags" => [
        "handler=#{@event['check']['handler']}",
      ],
      "text" => "#{@event['check']['output']}",
      "summary" => "#{action_to_string} - #{short_name}",
      "value" => "",
      "type" => "sensuAlert",
      "attributes" => {
        "subscribers" => "#{subscribers.join(",")}",
        "thresholdInfo" => "#{@event['action']}: #{@event['check']['command']}"
      },
      "rawData" => "#{@event.to_json}"
    }.to_json
    # puts payload

    headers = { 'Content-Type' => 'application/json' }
    if key
      headers['Authorization'] = 'Key ' + key
    end

    begin
      timeout 2 do
        response = HTTParty.post(url, :body => payload, :headers => headers)
        if response.success?
          puts 'alerta -- sent alert for ' + short_name + ' id: ' + response['id']
        else
          puts response
        end
      end
    rescue Timeout::Error
      puts 'alerta -- timed out while attempting to ' + @event['action'] + ' an incident -- ' + short_name
    end

  end
end
