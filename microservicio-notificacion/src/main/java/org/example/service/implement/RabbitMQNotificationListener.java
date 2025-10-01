package org.example.service.implement;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.rabbitmq.client.Channel;
import lombok.RequiredArgsConstructor;
import org.example.Config.RabbitMQConfig;
import org.example.dto.CreateNotificationDTO;
import org.example.model.Notification;
import org.example.service.interfaces.NotificationService;
import org.springframework.amqp.core.Message;
import org.springframework.amqp.rabbit.annotation.RabbitListener;
import org.springframework.stereotype.Service;

import jakarta.annotation.PostConstruct;

@Service
@RequiredArgsConstructor
public class RabbitMQNotificationListener {

    private final NotificationService notificationService;
    private final ObjectMapper objectMapper = new ObjectMapper();

    @PostConstruct
    public void init() {
        System.out.println("RabbitMQNotificationListener listo para recibir mensajes");
    }

    @RabbitListener(queues = RabbitMQConfig.NOTIFICATION_QUEUE, ackMode = "MANUAL")
    public void receiveNotification(Message message, Channel channel) throws Exception {
        try {
            // Convertimos el mensaje de JSON a DTO
            CreateNotificationDTO dto = objectMapper.readValue(message.getBody(), CreateNotificationDTO.class);

            // Creamos la notificación
            Notification notification = notificationService.createNotification(dto);

            // Enviamos a schedule
            notificationService.scheduleNotification(notification.getId());

            // Confirmamos a RabbitMQ que procesamos el mensaje correctamente
            channel.basicAck(message.getMessageProperties().getDeliveryTag(), false);

        } catch (Exception e) {
            // Rechazamos el mensaje para que no se quede bloqueado
            channel.basicReject(message.getMessageProperties().getDeliveryTag(), false);
            throw e; // opcional: loggear o manejar error
        }
    }
}
