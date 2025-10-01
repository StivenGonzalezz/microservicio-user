package org.example.service.implement;

import lombok.RequiredArgsConstructor;
import org.example.dto.CreateNotificationDTO;
import org.example.dto.EmailDTO;
import org.example.model.Notification;
import org.example.model.enums.Channels;
import org.example.model.enums.Status;
import org.example.repository.NotificationRepository;
import org.example.service.interfaces.EmailService;
import org.example.service.interfaces.NotificationService;
import org.example.service.interfaces.SMSService;
import org.example.service.interfaces.WhatsAppService;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Optional;

@Service
@RequiredArgsConstructor
public class NotificationServiceImpl implements NotificationService {

    private final NotificationRepository notificationRepository;
    private final EmailService emailService;
    private final SMSService smsService;
    private final WhatsAppService whatsAppService;

    @Override
    public Notification getNotificationFindById(String id) throws Exception {
     Notification notification = getById(id);
        return notification;
    }

    @Override
    public Page<Notification> getNotifications(int page, int size) {
        Pageable pageable = PageRequest.of(page, size, Sort.by("createdAt").descending());
        return notificationRepository.findAll(pageable);
    }

    @Override
    public List<Channels> getChannels(){

        return List.of(Channels.values());
    }

    @Override
    public Notification createNotification(CreateNotificationDTO createNotificationDTO)throws Exception{
        Notification notification = new Notification(

                createNotificationDTO.body(),
                Status.PENDING,
                LocalDateTime.now()
        );
        notification.setAffair(createNotificationDTO.affair());
        notification.setEmail(createNotificationDTO.email());
        notification.setNumber(createNotificationDTO.number());
        notification.setChannels("SMS,EMAIL,WHATSAPP");
        notificationRepository.save(notification);
        return notification;
    }
    @Override
    public String scheduleNotification(String id) throws Exception{

        Optional<Notification> notificationOptional = notificationRepository.findById(id);
        if (!notificationOptional.isEmpty()) {

            Notification notification = notificationOptional.get();
            if (notification.getStatus() != Status.SEND){
                String message = "";

                if (notification.getChannels().contains("EMAIL")) {
                    emailService.sendEmail(new EmailDTO(notification.getAffair(), notification.getMessage(), notification.getEmail()));
                    message += "sending email notification" + "\n";
                }
                if (notification.getChannels().contains("SMS")) {
                    smsService.sendSMS(notification.getNumber(), notification.getAffair() + "\n" + notification.getMessage());
                    message += "sending SMS notification" + "\n";
                }
                if (notification.getChannels().contains("WHATSAPP")) {
                    whatsAppService.sendMessage(notification.getNumber(), notification.getAffair() + "\n" + notification.getMessage());
                    message += "sending Whatsaap notification" + "\n";
                }
                if (message.equals("")) {
                    message += "communication channel not supported";
                } else {
                    notification.setSendAt(LocalDateTime.now());
                    notification.setStatus(Status.SEND);
                    notificationRepository.save(notification);
                }
                return message;
           }
            else {
                throw new Exception("the notification has already been sent");
            }
        }
        else {
            throw new  Exception ("The notification does no exits");
        }
    }

    public Notification getById(String id) throws Exception {

        Optional<Notification> notificationOptional = notificationRepository.findById(id);
        if (notificationOptional.isEmpty()){
            throw new Exception("The notification not exits");
        }
        return notificationOptional.get();
    }
}
