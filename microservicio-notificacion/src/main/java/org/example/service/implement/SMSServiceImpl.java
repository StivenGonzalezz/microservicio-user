package org.example.service.implement;

import com.twilio.Twilio;
import com.twilio.rest.api.v2010.account.Message;
import com.twilio.type.PhoneNumber;
import lombok.RequiredArgsConstructor;
import org.example.Config.TwilioConfig;
import org.example.service.interfaces.SMSService;
import org.springframework.stereotype.Service;

@Service
@RequiredArgsConstructor
public class SMSServiceImpl implements SMSService {

    private final TwilioConfig twilioConfig;
    private boolean initialized = false;

    private void initTwilio() {
        if (!initialized) {
            Twilio.init(twilioConfig.getAccountSid(), twilioConfig.getAuthToken());
            initialized = true;
        }
    }


    @Override
    public void sendSMS(String number, String message){
        initTwilio();

        Message.creator(
                new PhoneNumber(number),
                new PhoneNumber(twilioConfig.getFromNumber()),
                message
        ).create();

    }

}
