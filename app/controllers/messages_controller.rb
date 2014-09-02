class MessagesController < ApplicationController
  skip_before_filter :verify_authenticity_token

  before_filter :authorize, only: [:index]

  def index
    if current_user.email.in?(ENV['VALID_EMAILS'].split(','))
      @messages = Message.all.order('timestamp asc')
    else
      sign_out
      flash[:error] = 'Unauthorized.'
      redirect_to root_path
    end
  end

  def create
    message_params = params.symbolize_keys
    message_params.delete(:action)
    message_params.delete(:controller)

    Message.create!(message_params)

    render text: nil
  end

  def realtime
  end
end
