class MessagesController < ApplicationController
  def index
    Message.all.order('timestamp desc')
  end

  def create
    params.each do |k,v|
      puts "#{k}: #{v}"
    end
  end

  def realtime
  end
end
