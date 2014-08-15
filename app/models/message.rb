class Message < ActiveRecord::Base
  attr_accessor :token, :team_id, :channel_id, :channel_name, :timestamp,
                :user_id, :user_name, :text, :trigger_word
end
