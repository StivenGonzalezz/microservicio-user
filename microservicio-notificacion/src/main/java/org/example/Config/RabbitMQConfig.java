package org.example.Config;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.amqp.core.*;
import org.springframework.amqp.rabbit.connection.ConnectionFactory;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.amqp.support.converter.Jackson2JsonMessageConverter;

@Configuration
public class RabbitMQConfig {
    // --- Constantes ---
    public static final String NOTIFICATION_QUEUE = "messaging.user.notify";
    public static final String MESSAGING_EXCHANGE = "messaging.events";
    public static final String USER_EXCHANGE = "user.events";

    public static final String ROUTING_KEY_MESSAGING = "messaging.send";
    public static final String ROUTING_KEY_USER = "user.#";

    // --- Cola ---
    @Bean
    public Queue notificationQueue() {
        return QueueBuilder.durable(NOTIFICATION_QUEUE).build();
    }

    // --- Exchanges ---
    @Bean
    public TopicExchange messagingExchange() {
        return new TopicExchange(MESSAGING_EXCHANGE, true, false);
    }

    @Bean
    public TopicExchange userExchange() {
        return new TopicExchange(USER_EXCHANGE, true, false);
    }

    // --- Bindings ---
    @Bean
    public Binding bindingMessagingEvents(Queue notificationQueue, TopicExchange messagingExchange) {
        return BindingBuilder.bind(notificationQueue)
                .to(messagingExchange)
                .with(ROUTING_KEY_MESSAGING);
    }

    // --- Conversor JSON ---
    @Bean
    public Jackson2JsonMessageConverter jackson2JsonMessageConverter() {
        return new Jackson2JsonMessageConverter();
    }

    @Bean
    public RabbitTemplate rabbitTemplate(ConnectionFactory connectionFactory,
                                         Jackson2JsonMessageConverter messageConverter) {
        RabbitTemplate rabbitTemplate = new RabbitTemplate(connectionFactory);
        rabbitTemplate.setMessageConverter(messageConverter);
        return rabbitTemplate;
    }
}
